import { clsx, type ClassValue } from 'clsx'
import { twMerge } from 'tailwind-merge'
import ccTLDs from './cctlds';
import ms from 'ms';

export async function sleep(ms: number) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

interface DeviceInfo {
  type: "mobile" | "tablet" | "desktop" | "unknown";
  os: string;
  browser: string;
  isIOS: boolean;
  isAndroid: boolean;
  isMac: boolean;
  isWindows: boolean;
  isLinux: boolean;
  isSafari: boolean;
  isChrome: boolean;
  isFirefox: boolean;
  isEdge: boolean;
  model?: string;
}

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export const isValidUrl = (url: string) => {
  try {
    new URL(url);
    return true;
  } catch (e) {
    return false;
  }
};

export const getUrlFromString = (str: string) => {
  if (isValidUrl(str)) return str;
  try {
    if (str.includes(".") && !str.includes(" ")) {
      return new URL(`https://${str}`).toString();
    }
  } catch (_) {}
  return str;
};

export const getUrlFromStringIfValid = (str: string) => {
  if (isValidUrl(str)) return str;
  try {
    if (str.includes(".") && !str.includes(" ")) {
      return new URL(`https://${str}`).toString();
    }
  } catch (_) {}
  return null;
};

export const getSearchParams = (url: string) => {
  // Create a params object
  const params = {} as Record<string, string>;

  new URL(url).searchParams.forEach((val, key) => {
    params[key] = val;
  });

  return params;
};

export const getSearchParamsWithArray = (url: string) => {
  const params = {} as Record<string, string | string[]>;

  new URL(url).searchParams.forEach((val, key) => {
    if (key in params) {
      const param = params[key];
      Array.isArray(param) ? param.push(val) : (params[key] = [param, val]);
    } else {
      params[key] = val;
    }
  });

  return params;
};

export const getPrettyUrl = (url?: string | null) => {
  if (!url) return "";
  // remove protocol (http/https) and www.
  // also remove trailing slash
  return url
    .replace(/(^\w+:|^)\/\//, "")
    .replace("www.", "")
    .replace(/\/$/, "");
};

export function inferPlatform(url: string): string {
  if (!url) return "Direct";

  try {
    // Remove protocol and www, get hostname
    let hostname = url
      .toLowerCase()
      .replace(/^(https?:\/\/)?(www\.)?/, "")
      .split("/")[0];

    // Remove port if exists
    hostname = hostname.split(":")[0];

    // Special cases mapping
    const platformMap: { [key: string]: string } = {
      "google.com": "Google",
      "facebook.com": "Facebook",
      "instagram.com": "Instagram",
      "twitter.com": "Twitter",
      "x.com": "Twitter",
      "linkedin.com": "LinkedIn",
      "youtube.com": "YouTube",
      "youtu.be": "YouTube",
      "pinterest.com": "Pinterest",
      "reddit.com": "Reddit",
      "tiktok.com": "TikTok",
      "github.com": "GitHub",
      "medium.com": "Medium",
      "dev.to": "Dev.to",
      "t.co": "Twitter",
      "fb.com": "Facebook",
      "fb.me": "Facebook",
      "bloomberg.com": "Bloomberg",
      "reuters.com": "Reuters",
      "nytimes.com": "New York Times",
      "wsj.com": "Wall Street Journal",
      "duckduckgo.com": "DuckDuckGo",
      "bing.com": "Bing",
      "yahoo.com": "Yahoo",
      "mail.google.com": "Gmail",
      "outlook.com": "Outlook",
      "substack.com": "Substack",
    };

    // Check for exact matches first
    if (platformMap[hostname]) {
      return platformMap[hostname];
    }

    // Check for subdomains of known platforms
    for (const [domain, platform] of Object.entries(platformMap)) {
      if (hostname.endsWith(`.${domain}`)) {
        return platform;
      }
    }

    // For unknown domains, capitalize first letter of each word
    // Remove common TLDs first
    hostname =
      hostname
        .replace(/\.(com|org|net|edu|gov|io|co|me|app|dev)$/, "")
        .split(".")
        .pop() || hostname;

    return hostname
      .split(/[-_]/)
      .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
      .join(" ");
  } catch (error) {
    console.error("Error inferring platform:", error);
    return "Unknown";
  }
}

export function inferSlug(url: string): string | null {
  try {
    const pathname = new URL(url).pathname;
    const slug = pathname.replace(/^\/|\/$/g, ""); // Remove leading and trailing slashes
    return slug || null; // Return null if the slug is empty
  } catch (error) {
    console.error("Invalid URL:", error);
    return null;
  }
}

export const getTldFromUrl = (url: string): string | null => {
  try {
    const parsedUrl = new URL(url);
    const hostname = parsedUrl.hostname;
    const tld = hostname.substring(hostname.lastIndexOf(".") + 1);

    return tld.length > 0 ? tld : null;
  } catch (error) {
    console.error("Invalid URL", error);
    return null;
  }
};

export function detectDevice(): DeviceInfo {
  const userAgent = navigator.userAgent.toLowerCase();

  // Initialize result with defaults
  const result: DeviceInfo = {
    type: "unknown",
    os: "unknown",
    browser: "unknown",
    isIOS: false,
    isAndroid: false,
    isMac: false,
    isWindows: false,
    isLinux: false,
    isSafari: false,
    isChrome: false,
    isFirefox: false,
    isEdge: false,
  };

  // Device type detection
  if (
    /(ipad|tablet|(android(?!.*mobile))|(windows(?!.*phone)(.*touch))|kindle|playbook|silk|tablet|(puffin(?!.*(IP|AP|WP))))/.test(
      userAgent,
    )
  ) {
    result.type = "tablet";
  } else if (
    /(mobi|ipod|phone|blackberry|opera mini|fennec|minimo|symbian|psp|nintendo ds|archos|skyfire|puffin|blazer|bolt|gobrowser|iris|maemo|semc|teashark|uzard)/.test(
      userAgent,
    )
  ) {
    result.type = "mobile";
  } else {
    result.type = "desktop";
  }

  // OS detection
  if (/iphone|ipad|ipod/.test(userAgent)) {
    result.os = "iOS";
    result.isIOS = true;
  } else if (/android/.test(userAgent)) {
    result.os = "Android";
    result.isAndroid = true;
  } else if (/macintosh|mac os x/.test(userAgent)) {
    result.os = "macOS";
    result.isMac = true;
  } else if (/windows|win32|win64|wow64/.test(userAgent)) {
    result.os = "Windows";
    result.isWindows = true;
  } else if (/linux/.test(userAgent) && !result.isAndroid) {
    result.os = "Linux";
    result.isLinux = true;
  }

  // Browser detection
  if (/edg/.test(userAgent)) {
    result.browser = "Edge";
    result.isEdge = true;
  } else if (/chrome/.test(userAgent) && !/chromium|edg/.test(userAgent)) {
    result.browser = "Chrome";
    result.isChrome = true;
  } else if (/firefox/.test(userAgent)) {
    result.browser = "Firefox";
    result.isFirefox = true;
  } else if (
    /safari/.test(userAgent) &&
    !/chrome|chromium|edg/.test(userAgent)
  ) {
    result.browser = "Safari";
    result.isSafari = true;
  } else if (/msie|trident/.test(userAgent)) {
    result.browser = "Internet Explorer";
  } else if (/opera/.test(userAgent)) {
    result.browser = "Opera";
  }

  // Try to detect models for common devices
  if (result.isIOS) {
    const matches =
      userAgent.match(/iphone\s+os\s+(\d+)_(\d+)/i) ||
      userAgent.match(/ipad;\s+cpu\s+os\s+(\d+)_(\d+)/i);
    if (matches) {
      const model = userAgent.includes("ipad") ? "iPad" : "iPhone";
      result.model = `${model} (iOS ${matches[1]}.${matches[2]})`;
    }
  } else if (result.isAndroid) {
    const matches = userAgent.match(/android\s+(\d+)(\.(\d+))?/i);
    if (matches) {
      result.model = `Android ${matches[1]}${matches[3] ? `.${matches[3]}` : ""}`;
    }
  }

  return result;
}



export const SECOND_LEVEL_DOMAINS = new Set([
  "com",
  "co",
  "net",
  "org",
  "edu",
  "gov",
  "in",
]);

export const SPECIAL_APEX_DOMAINS = {
  "youtu.be": "youtube.com",
};

export const GOOGLE_FAVICON_URL =
  "https://www.google.com/s2/favicons?sz=64&domain_url=";


export const getApexDomain = (url: string) => {
  let domain: any;
  try {
    domain = new URL(url).hostname;
  } catch (e) {
    return "";
  }
  // special apex domains (e.g. youtu.be)
  if (SPECIAL_APEX_DOMAINS[domain as keyof typeof SPECIAL_APEX_DOMAINS])
    return SPECIAL_APEX_DOMAINS[domain as keyof typeof SPECIAL_APEX_DOMAINS];

  const parts = domain.split(".");
  if (parts.length > 2) {
    // if this is a second-level TLD (e.g. co.uk, .com.ua, .org.tt), we need to return the last 3 parts
    if (
      SECOND_LEVEL_DOMAINS.has(parts[parts.length - 2]) &&
      ccTLDs.has(parts[parts.length - 1])
    ) {
      return parts.slice(-3).join(".");
    }
    // otherwise, it's a subdomain (e.g. dub.vercel.app), so we return the last 2 parts
    return parts.slice(-2).join(".");
  }
  // if it's a normal domain (e.g. dub.sh), we return the domain
  return domain;
};

// Verify the URL entered by user
export const validDomainRegex = new RegExp(
  /^([a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}$/,
);

// Get time link was added
export const timeAgo = (timestamp: string) => {
  if (!timestamp) return "Just now";
  const diff = Date.now() - new Date(timestamp).getTime();
  if (diff < 60000) {
    // less than 1 second
    return "Just now";
  }
  if (diff > 82800000) {
    // more than 23 hours – similar to how Twitter displays timestamps
    return new Date(timestamp).toLocaleDateString("en-US", {
      month: "short",
      day: "numeric",
      year:
        new Date(timestamp).getFullYear() !== new Date().getFullYear()
          ? "numeric"
          : undefined,
    });
  }
  return ms(diff);
};

export const removeHashFromHexColor = (hexColor: string) => {
  // Use a regular expression to match the # symbol at the beginning
  return hexColor.replace(/^#/, "");
};

export const getCurrentBaseURL = () => {
  if (typeof window !== "undefined") {
    const baseURL = window.location.origin;
    return baseURL;
  }
};

export const getInitials = (name: string | undefined): string => {
  if (!name) return "";

  const words = name.split(" ");
  if (words.length > 1) {
    return (words[0][0] + words[1][0]).toUpperCase();
  }
  return name.slice(0, 2).toUpperCase();
};




