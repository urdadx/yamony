import { argon2id } from "@noble/hashes/argon2.js";
import { randomBytes as nobleRandomBytes } from "@noble/hashes/utils.js";

/**
 * Generate random bytes
 */
function randomBytes(length: number): Uint8Array {
  return nobleRandomBytes(length);
}

/**
 * Convert Uint8Array to ArrayBuffer for Web Crypto API
 */
function toArrayBuffer(arr: Uint8Array): ArrayBuffer {
  // Create a new ArrayBuffer copy to avoid SharedArrayBuffer type issues
  const buffer = new ArrayBuffer(arr.byteLength);
  new Uint8Array(buffer).set(arr);
  return buffer;
}

/**
 * Generate a random 256-bit Vault Encryption Key (VEK)
 */
export function generateVEK(): Uint8Array {
  return randomBytes(32); // 256 bits
}

/**
 * Derive a master key from the user's password using Argon2id
 */
export async function deriveMasterKey(
  password: string,
  salt: Uint8Array
): Promise<Uint8Array> {
  const encoder = new TextEncoder();
  const passwordBytes = encoder.encode(password);

  // Argon2id parameters (matching backend)
  const masterKey = argon2id(passwordBytes, salt, {
    t: 3, // iterations (time cost)
    m: 65536, // memory cost in KiB (64 MB)
    p: 2, // parallelism
  });

  return masterKey;
}

/**
 * Generate a random salt for Argon2id
 */
export function generateSalt(): Uint8Array {
  return randomBytes(32); // 256 bits
}

/**
 * Wrap (encrypt) the VEK with the master key using AES-256-GCM
 */
export async function wrapVEK(
  vek: Uint8Array,
  masterKey: Uint8Array
): Promise<{
  wrappedVEK: Uint8Array;
  iv: Uint8Array;
  tag: Uint8Array;
}> {
  const iv = randomBytes(12); // 96 bits for GCM

  // Import master key for AES-GCM
  const cryptoKey = await crypto.subtle.importKey(
    "raw",
    toArrayBuffer(masterKey),
    { name: "AES-GCM" },
    false,
    ["encrypt"]
  );

  // Encrypt VEK
  const encrypted = await crypto.subtle.encrypt(
    {
      name: "AES-GCM",
      iv: toArrayBuffer(iv),
      tagLength: 128, // 16 bytes
    },
    cryptoKey,
    toArrayBuffer(vek)
  );

  // Split encrypted data into ciphertext and auth tag
  const encryptedArray = new Uint8Array(encrypted);
  const wrappedVEK = encryptedArray.slice(0, -16); // Everything except last 16 bytes
  const tag = encryptedArray.slice(-16); // Last 16 bytes

  return { wrappedVEK, iv, tag };
}

/**
 * Unwrap (decrypt) the VEK with the master key using AES-256-GCM
 */
export async function unwrapVEK(
  wrappedVEK: Uint8Array,
  masterKey: Uint8Array,
  iv: Uint8Array,
  tag: Uint8Array
): Promise<Uint8Array> {
  // Import master key for AES-GCM
  const cryptoKey = await crypto.subtle.importKey(
    "raw",
    toArrayBuffer(masterKey),
    { name: "AES-GCM" },
    false,
    ["decrypt"]
  );

  // Combine ciphertext and tag
  const combined = new Uint8Array(wrappedVEK.length + tag.length);
  combined.set(wrappedVEK);
  combined.set(tag, wrappedVEK.length);

  // Decrypt VEK
  const decrypted = await crypto.subtle.decrypt(
    {
      name: "AES-GCM",
      iv: toArrayBuffer(iv),
      tagLength: 128,
    },
    cryptoKey,
    toArrayBuffer(combined)
  );

  return new Uint8Array(decrypted);
}

/**
 * Convert Uint8Array to base64 string
 */
export function toBase64(arr: Uint8Array): string {
  return btoa(String.fromCharCode(...arr));
}

/**
 * Convert base64 string to Uint8Array
 */
export function fromBase64(str: string): Uint8Array {
  return Uint8Array.from(atob(str), (c) => c.charCodeAt(0));
}
