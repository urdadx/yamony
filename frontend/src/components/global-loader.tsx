import Spinner from "./ui/spinner";

export function GlobalLoader() {
  return (
    <div className="w-full p-2 sm:p-4">

      <div className="fixed inset-0 flex items-center justify-center">
        <Spinner size={28} className="text-gray-400" />
      </div>
    </div>
  );
}