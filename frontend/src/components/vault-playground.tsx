import { FilterByCategories } from "./filter-by-categories";
import { GlobalSearch } from "./global-search";
import { VaultItem } from "./vault-item";

export const VaultPlayground = () => {

  return (
    <div className="w-[450px] ">
      <div className="flex p-2 gap-2 items-center justify-between">
        <GlobalSearch />
        <FilterByCategories />
        {/* <SortFilter /> */}
      </div>
      <div className="py-2">
        <VaultItem />
        <VaultItem />
        <VaultItem />
      </div>

    </div>
  );
};

