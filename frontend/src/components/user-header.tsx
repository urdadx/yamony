import { Plus, Youtube } from 'lucide-react'
import { Button } from './ui/button'
import { AddIcon } from '@/assets/icons/add-icon'
import { Avatar, AvatarFallback, AvatarImage } from './ui/avatar'
import { CameraIcon } from '@/assets/icons/camera-icon'

export function HomeHeader() {
  return (
    <div className="h-fit flex items-center justify-center ">
      <div className="w-full">
        <div className="flex items-center gap-4 mb-6">
          <Avatar className="w-20 h-20 boder-4 border-white">
            <AvatarImage src="/path-to-avatar-image.jpg" alt="Rogue Shinobi" />
            <AvatarFallback className="bg-gray-100">
              <CameraIcon color='#84cc16' className='w-6 h-6 text-primary' />
            </AvatarFallback>
          </Avatar>
          <div className="flex-1">
            <h1 className="text-xl font-bold text-gray-900 mb-3">Rogue Shinobi</h1>
            <div className="flex items-center gap-3">
              <div className="relative">
                <Button size="icon" className="w-10 h-10 rounded-full bg-gray-200 hover:bg-gray-300 transition flex items-center justify-center">
                  <Youtube className="w-5 h-5 text-gray-600" />
                </Button>

              </div>

              <Button size="icon" className="w-10 h-10 rounded-full bg-gray-200 hover:bg-gray-300 transition flex items-center justify-center">
                <Plus className="w-5 h-5 text-gray-600" />
              </Button>
            </div>
          </div>
        </div>

        <Button className="w-full rounded-3xl p-5">
          <AddIcon className='w-10 h-10 size-6!' color="#FFFFFF" />
          <span className='text-md'>
            Add block
          </span>
        </Button>
      </div>
    </div>
  )
}
