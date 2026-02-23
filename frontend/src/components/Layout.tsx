import React from 'react'
import { Link, useLocation, useNavigate } from 'react-router-dom'
import { CalendarDaysIcon, BellIcon, ArrowRightOnRectangleIcon, PlusIcon } from '@heroicons/react/24/outline'
import { useAuthStore } from '../store/authStore'
import { logout } from '../api/auth'
import { useIncomingInvitations } from '../hooks/useInvitations'

interface LayoutProps {
  children: React.ReactNode
}

export const Layout: React.FC<LayoutProps> = ({ children }) => {
  const { user, logout: storeLogout } = useAuthStore()
  const navigate = useNavigate()
  const location = useLocation()
  const { data: incoming } = useIncomingInvitations()
  const pendingCount = incoming?.invitations?.length ?? 0

  const handleLogout = async () => {
    try { await logout() } catch { /* ignore */ }
    storeLogout()
    navigate('/login')
  }

  const navLink = (to: string, label: string, icon: React.ReactNode) => {
    const active = location.pathname.startsWith(to)
    return (
      <Link
        to={to}
        className={`flex items-center gap-2 px-3 py-2 rounded-lg text-sm font-medium transition-colors ${
          active ? 'bg-blue-50 text-blue-700' : 'text-gray-600 hover:bg-gray-100'
        }`}
      >
        {icon}
        {label}
      </Link>
    )
  }

  return (
    <div className="min-h-screen bg-gray-50 flex">
      {/* Sidebar */}
      <aside className="w-60 bg-white border-r flex flex-col py-4 px-3 gap-1 fixed h-full">
        <div className="px-3 mb-4">
          <h1 className="text-xl font-bold text-blue-600">EventScheduler</h1>
          <p className="text-xs text-gray-500 mt-0.5">{user?.username}</p>
        </div>

        <Link
          to="/events/new"
          className="flex items-center gap-2 px-3 py-2 bg-blue-600 text-white rounded-lg text-sm font-medium hover:bg-blue-700 mb-2"
        >
          <PlusIcon className="h-4 w-4" />
          New Event
        </Link>

        {navLink('/events', 'Events', <CalendarDaysIcon className="h-4 w-4" />)}

        <Link
          to="/invitations"
          className={`flex items-center gap-2 px-3 py-2 rounded-lg text-sm font-medium transition-colors ${
            location.pathname.startsWith('/invitations')
              ? 'bg-blue-50 text-blue-700'
              : 'text-gray-600 hover:bg-gray-100'
          }`}
        >
          <div className="relative">
            <BellIcon className="h-4 w-4" />
            {pendingCount > 0 && (
              <span className="absolute -top-1 -right-1 h-3 w-3 bg-red-500 rounded-full text-white text-[9px] flex items-center justify-center">
                {pendingCount > 9 ? '9+' : pendingCount}
              </span>
            )}
          </div>
          Invitations
          {pendingCount > 0 && (
            <span className="ml-auto bg-red-100 text-red-700 text-xs rounded-full px-1.5 py-0.5">
              {pendingCount}
            </span>
          )}
        </Link>

        <div className="mt-auto">
          <button
            onClick={handleLogout}
            className="flex items-center gap-2 w-full px-3 py-2 rounded-lg text-sm font-medium text-gray-600 hover:bg-gray-100"
          >
            <ArrowRightOnRectangleIcon className="h-4 w-4" />
            Sign out
          </button>
        </div>
      </aside>

      {/* Main content */}
      <main className="ml-60 flex-1 p-6">
        {children}
      </main>
    </div>
  )
}
