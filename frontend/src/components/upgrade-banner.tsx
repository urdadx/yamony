import { SolarStarIcon } from "@/assets/icons/star-icon";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { TextShimmer } from "./text-shimmer";

export function UpgradeBanner() {

  return (
    <Card className="shadow-none h-fit gap-2">
      <div className="px-4 ">
        <CardTitle className="text-sm flex gap-1 font-normal">
          <SolarStarIcon size="20" color="#84cc16" />
          Upgrade to Pro
        </CardTitle>
        <CardDescription className="text-sm py-2 leading-tight text-muted-foreground">
          Upgrade now to continue enjoying all features.
        </CardDescription>
      </div>
      <div className="grid px-3">
        <a target="_blank" href="/choose-plan" rel="noreferrer">
          <Button
            className="w-full text-sidebar-primary-foreground shadow-none"
            size="sm"
          >
            <TextShimmer duration={1.5} className="text-white">
              Upgrade your plan
            </TextShimmer>
          </Button>
        </a>
      </div>
    </Card>
  );
}