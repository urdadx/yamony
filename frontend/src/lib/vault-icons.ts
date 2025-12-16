import { HomeIcon } from "@/assets/icons/home-icon";
import { ShieldStrongIcon } from "@/assets/icons/shield-strong-icon";
import { SolarStarIcon } from "@/assets/icons/star-icon";
import { CameraIcon } from "@/assets/icons/camera-icon";
import { BoltIcon } from "@/assets/icons/bolt-icon";
import { PaletteIcon } from "@/assets/icons/palette-icon";
import { PhoneIcon } from "@/assets/icons/phone-icon";
import { LetterIcon } from "@/assets/icons/letter";
import { SolarBuildingsBoldDuotone } from "@/assets/icons/building-icon";
import { SolarEarthBoldDuotone } from "@/assets/icons/globe";

export const VAULT_ICONS = {
  Home: HomeIcon,
  Shield: ShieldStrongIcon,
  Star: SolarStarIcon,
  Earth: SolarEarthBoldDuotone,
  Buildings: SolarBuildingsBoldDuotone,
  Bolt: BoltIcon,
  Palette: PaletteIcon,
  Phone: PhoneIcon,
  Letter: LetterIcon,
  Camera: CameraIcon,
  
};

export type VaultIconName = keyof typeof VAULT_ICONS;
