import React from 'react'
import { Link } from 'react-router-dom'
import { format } from 'date-fns'
import { MapPinIcon, ClockIcon, UsersIcon } from '@heroicons/react/24/outline'
import type { Event } from '../types'
import { StatusBadge } from './StatusBadge'
import { useAuthStore } from '../store/authStore'

interface EventCardProps {
  event: Event
  onDelete?: (id: string) => void
}

export const EventCard: React.FC<EventCardProps> = ({ event, onDelete }) => {
  const userId = useAuthStore((s) => s.user?.id)
  const isOwner = event.owner_id === userId
  const myAttendee = event.attendees?.find((a) => a.user_id === userId)
  const status = myAttendee?.status ?? 'upcoming'

  return (
    <div className="bg-white rounded-xl border shadow-sm hover:shadow-md transition-shadow p-4">
      <div className="flex items-start justify-between gap-2 mb-2">
        <div className="flex items-center gap-2 min-w-0">
          <div
            className="w-3 h-3 rounded-full flex-shrink-0"
            style={{ backgroundColor: event.color }}
          />
          <Link
            to={`/events/${event.id}`}
            className="font-semibold text-gray-900 hover:text-blue-600 truncate"
          >
            {event.title}
          </Link>
        </div>
        <StatusBadge status={status} />
      </div>

      <div className="space-y-1 text-sm text-gray-500">
        <div className="flex items-center gap-1.5">
          <ClockIcon className="h-4 w-4 flex-shrink-0" />
          <span>
            {event.is_all_day
              ? `All day · ${format(new Date(event.starts_at), 'MMM d, yyyy')}`
              : `${format(new Date(event.starts_at), 'MMM d, h:mm a')} – ${format(new Date(event.ends_at), 'h:mm a')}`}
          </span>
        </div>
        {event.location && (
          <div className="flex items-center gap-1.5">
            <MapPinIcon className="h-4 w-4 flex-shrink-0" />
            <span className="truncate">{event.location}</span>
          </div>
        )}
        {event.attendees && event.attendees.length > 0 && (
          <div className="flex items-center gap-1.5">
            <UsersIcon className="h-4 w-4 flex-shrink-0" />
            <span>{event.attendees.length} attendee{event.attendees.length !== 1 ? 's' : ''}</span>
          </div>
        )}
      </div>

      <div className="flex items-center gap-2 mt-3 pt-3 border-t">
        <Link
          to={`/events/${event.id}`}
          className="text-xs text-blue-600 hover:underline"
        >
          View details
        </Link>
        {isOwner && (
          <>
            <Link
              to={`/events/${event.id}/edit`}
              className="text-xs text-gray-500 hover:underline"
            >
              Edit
            </Link>
            <button
              onClick={() => onDelete?.(event.id)}
              className="text-xs text-red-500 hover:underline ml-auto"
            >
              Delete
            </button>
          </>
        )}
      </div>
    </div>
  )
}
