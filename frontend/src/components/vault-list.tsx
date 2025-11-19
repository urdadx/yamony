import { FilterByCategories } from "./filter-by-categories";
import { VaultItem } from "./vault-item";

export const VaultList = () => {

  return (
    <div className="w-[370px] ">
      <FilterByCategories />
      <div className="py-2">
        <VaultItem />
        <VaultItem />
        <VaultItem />
      </div>

    </div>
  );
};

