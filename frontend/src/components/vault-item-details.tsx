import { MoreVertical, Eye, EyeOff, Mail, Lock, Globe2Icon, Wand2Icon, PencilIcon, ZapIcon, StickyNote } from "lucide-react"
import { Button } from "@/components/ui/button"
import { useState } from "react";
import { ShieldWeakIcon } from "@/assets/icons/shield-weak-icon";
import { UserPlusIcon } from "@/assets/icons/user-plus-icon";
import { EditPencilIcon } from "@/assets/icons/edit-pencil-icon";

export function VaultItemDetails() {
  const [showPassword, setShowPassword] = useState(false);
  return (
    <main className="min-h-screen w-full p-4">
      <div className=" w-full">
        <div className="flex items-center justify-between mb-6">
          <h1 className="text-lg font-bold text-gray-900">wd1.myworkdayjobs.com</h1>
          <div className="flex items-center gap-2">
            <Button variant="outline" className=" text-primary rounded-lg shadow-none p-4 hover:bg-rose-50/40 hover:text-primary/80">
              <EditPencilIcon color="#f43f5e" className="hover:scale-110 transition-transform" />
              <span className="hover:opacity-80 font-normal transition-opacity">Edit</span>
            </Button>
            <Button variant="outline" className="text-primary rounded-lg shadow-none p-3 hover:bg-rose-50/40 hover:text-primary/80">
              <UserPlusIcon color="#f43f5e" className="hover:scale-110 transition-transform" />
            </Button>
            <Button variant="outline" className="text-primaryrounded-lg shadow-none p-3 hover:bg-rose-50/40 hover:text-primary/80">
              <MoreVertical size={12} className="hover:scale-110 text-primary transition-transform" />
            </Button>
          </div>
        </div>

        <div className=" w-full overflow-hidden bg-white border rounded-2xl">

          <div className="flex items-center gap-4 px-4 py-4 border-b hover:bg-gray-50 transition-colors cursor-pointer">
            <Mail className="h-4 w-4 text-gray-500 shrink-0" />

            <div className="flex-1">
              <label className="block text-xs font-medium text-gray-500 mb-1">
                Email
              </label>
              <div className="text-sm text-gray-900">
                abassabdulwahab3@gmail.com
              </div>
            </div>
          </div>

          <div className="flex items-center gap-4 px-4 py-3 hover:bg-gray-50 transition-colors cursor-pointer">
            <Lock className="h-4 w-4 text-gray-500 shrink-0" />

            <div className="flex-1">
              <label className="block text-xs font-medium text-gray-500">
                Password
              </label>

              <div className="flex items-center justify-between">
                <div className="text-sm text-gray-900 flex-1">
                  {showPassword ? 'MyPassword123!' : '••••••••••••'}
                </div>
                <div className="flex items-center gap-2">
                  <div className="flex items-center gap-1">
                    <ShieldWeakIcon color="#f97316" className="h-5 w-5" />
                    <span className="text-sm text-orange-400 font-medium">Weak password</span>
                  </div>

                  <Button
                    variant="ghost"
                    size="icon"
                    className="h-8 w-8"
                    onClick={() => setShowPassword(!showPassword)}
                  >
                    {showPassword ? (
                      <EyeOff className="h-5 w-5" />
                    ) : (
                      <Eye className="h-5 w-5" />
                    )}
                  </Button>
                </div>
              </div>
            </div>
          </div>

        </div>
        <div className="flex mt-3 bg-white rounded-2xl border items-center gap-4 px-4 py-4 hover:bg-gray-50 transition-colors cursor-pointer">
          <Globe2Icon className="h-4 w-4 text-gray-500 shrink-0" />

          <div className="flex-1">
            <label className="block text-xs font-medium text-gray-500 mb-1">
              Websites
            </label>
            <div className="text-sm text-gray-900">
              <a href="https://wd1.myworkdayjobs.com/en-US/WD1" className="text-indigo-600 hover:underline">https://wd1.myworkdayjobs.com/en-US/WD1</a>
            </div>
          </div>
        </div>
        <div className="flex mt-3 bg-white rounded-2xl border items-center gap-4 px-4 py-4 hover:bg-gray-50 transition-colors cursor-pointer">
          <StickyNote className="h-4 w-4 text-gray-500 shrink-0" />

          <div className="flex-1">
            <label className="block text-xs font-medium text-gray-500 mb-1">
              Notes
            </label>
            <div className="text-sm text-gray-900">
              This is my super secret details
            </div>
          </div>
        </div>

        <div className=" mt-3 w-full overflow-hidden border rounded-2xl">

          <div className="flex items-center gap-4 px-4 py-3 hover:bg-gray-50 transition-colors cursor-pointer">
            <Wand2Icon className="h-4 w-4 text-gray-500 shrink-0" />

            <div className="flex-1">
              <label className="block text-xs font-medium text-gray-500 mb-1">
                Last autofill
              </label>
              <div className="text-sm text-gray-900">
                Nov 23, 2025
              </div>
            </div>
          </div>
          <div className="flex items-center gap-4 px-4 py-3 hover:bg-gray-50 transition-colors cursor-pointer">
            <PencilIcon className="h-4 w-4 text-gray-500 shrink-0" />

            <div className="flex-1">
              <label className="block text-xs font-medium text-gray-500 mb-1">
                Last edited
              </label>
              <div className="text-sm text-gray-900">
                Sept 23, 2025
              </div>
            </div>
          </div>
          <div className="flex items-center gap-4 px-4 py-3 hover:bg-gray-50 transition-colors cursor-pointer">
            <ZapIcon className="h-4 w-4 text-gray-500 shrink-0" />

            <div className="flex-1">
              <label className="block text-xs font-medium text-gray-500 mb-1">
                Created on
              </label>
              <div className="text-sm text-gray-900">
                Jan 23, 2025
              </div>
            </div>
          </div>


        </div>

      </div>
    </main>
  )
}


