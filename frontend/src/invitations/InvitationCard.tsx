import React from 'react'
import { format } from 'date-fns'
import { CalendarDaysIcon, MapPinIcon, EnvelopeIcon } from '@heroicons/react/24/outline'
import type { Invitation } from '../types'
import { Button } from '../components/Button'
import { useRespondInvitation } from '../hooks/useInvitations'

interface InvitationCardProps {
  invitation: Invitation
}

const statusColors: Record<string, string> = {
  pending: 'bg-amber-100 text-amber-700',
  accepted: 'bg-green-100 text-green-700',
  declined: 'bg-gray-100 text-gray-500',
  expired: 'bg-red-100 text-red-500',
}

export const InvitationCard: React.FC<InvitationCardProps> = ({ invitation }) => {
  const respond = useRespondInvitation()
  const event = invitation.event

  const handleRespond = (accept: boolean) => {
    if (!invitation.token) return
    respond.mutate({ token: invitation.token, accept })
  }

  return (
    <div className="bg-white rounded-xl border shadow-sm p-4 space-y-3">
      <div className="flex items-start justify-between gap-2">
        <div>
          <h3 className="font-semibold text-gray-900">{event?.title ?? 'Event invitation'}</h3>
          {invitation.inviter && (
            <p className="text-xs text-gray-500">from {invitation.inviter.username ?? invitation.inviter_id}</p>
          )}
        </div>
        <span className={`text-xs font-medium px-2 py-0.5 rounded-full ${statusColors[invitation.status] ?? ''}`}>
          {invitation.status}
        </span>
      </div>

      {event && (
        <div className="space-y-1 text-sm text-gray-600">
          <div className="flex items-center gap-2">
            <CalendarDaysIcon className="h-4 w-4 text-gray-400" />
            <span>
              {event.is_all_day
                ? format(new Date(event.starts_at), 'MMM d, yyyy')
                : `${format(new Date(event.starts_at), 'MMM d, h:mm a')} – ${format(new Date(event.ends_at), 'h:mm a')}`}
            </span>
          </div>
          {event.location && (
            <div className="flex items-center gap-2">
              <MapPinIcon className="h-4 w-4 text-gray-400" />
              <span>{event.location}</span>
            </div>
          )}
        </div>
      )}

      {invitation.message && (
        <div className="flex items-start gap-2 text-sm text-gray-600 bg-gray-50 rounded-lg p-2">
          <EnvelopeIcon className="h-4 w-4 text-gray-400 flex-shrink-0 mt-0.5" />
          <span>{invitation.message}</span>
        </div>
      )}

      {invitation.status === 'pending' && (
        <div className="flex gap-2 pt-1">
          <Button
            onClick={() => handleRespond(true)}
            loading={respond.isPending}
            size="sm"
          >
            Accept
          </Button>
          <Button
            variant="secondary"
            onClick={() => handleRespond(false)}
            loading={respond.isPending}
            size="sm"
          >
            Decline
          </Button>
        </div>
      )}

      <p className="text-xs text-gray-400">
        Expires {format(new Date(invitation.expires_at), 'MMM d, yyyy')}
      </p>
    </div>
  )
}
