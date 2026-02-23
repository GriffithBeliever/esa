import React, { useState } from 'react'
import { useParams, useNavigate, Link } from 'react-router-dom'
import { useQuery, useQueryClient } from '@tanstack/react-query'
import { format } from 'date-fns'
import { getEvent } from '../api/events'
import { useDeleteEvent, useUpdateEventStatus } from '../hooks/useEvents'
import { useAuthStore } from '../store/authStore'
import { Button } from '../components/Button'
import { Spinner } from '../components/Spinner'
import { StatusBadge } from './StatusBadge'
import { MapPinIcon, ClockIcon, UsersIcon, PencilSquareIcon, TrashIcon } from '@heroicons/react/24/outline'
import { InvitationList } from '../invitations/InvitationList'
import { InviteUserModal } from '../invitations/InviteUserModal'
import type { EventStatus } from '../types'

export const EventDetail: React.FC = () => {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const userId = useAuthStore((s) => s.user?.id)
  const qc = useQueryClient()

  const { data: event, isLoading } = useQuery({
    queryKey: ['event', id],
    queryFn: () => getEvent(id!),
    enabled: !!id,
  })

  const deleteEvent = useDeleteEvent()
  const updateStatus = useUpdateEventStatus()
  const [showInviteModal, setShowInviteModal] = useState(false)

  if (isLoading) return <div className="flex justify-center p-12"><Spinner /></div>
  if (!event) return <div className="text-center p-12 text-gray-500">Event not found</div>

  const isOwner = event.owner_id === userId
  const myAttendee = event.attendees?.find((a) => a.user_id === userId)
  const myStatus = myAttendee?.status ?? 'upcoming'

  const handleDelete = async () => {
    if (!confirm('Delete this event? This cannot be undone.')) return
    await deleteEvent.mutateAsync(event.id)
    navigate('/events')
  }

  const handleStatusChange = (status: EventStatus) => {
    updateStatus.mutate({ id: event.id, status }, {
      onSuccess: () => qc.invalidateQueries({ queryKey: ['event', id] }),
    })
  }

  const STATUS_OPTIONS: { value: EventStatus; label: string }[] = [
    { value: 'attending', label: 'Attending' },
    { value: 'maybe', label: 'Maybe' },
    { value: 'declined', label: 'Not going' },
  ]

  return (
    <div className="max-w-2xl mx-auto space-y-6">
      {/* Header */}
      <div className="flex items-start justify-between gap-4">
        <div className="flex items-center gap-3">
          <div className="w-4 h-4 rounded-full flex-shrink-0 mt-1" style={{ backgroundColor: event.color }} />
          <div>
            <h1 className="text-2xl font-bold text-gray-900">{event.title}</h1>
            {event.owner && (
              <p className="text-sm text-gray-500">by {event.owner.username}</p>
            )}
          </div>
        </div>
        <div className="flex items-center gap-2">
          {isOwner && (
            <>
              <Link to={`/events/${event.id}/edit`}>
                <Button variant="secondary" size="sm">
                  <PencilSquareIcon className="h-4 w-4 mr-1" />
                  Edit
                </Button>
              </Link>
              <Button variant="danger" size="sm" onClick={handleDelete} loading={deleteEvent.isPending}>
                <TrashIcon className="h-4 w-4 mr-1" />
                Delete
              </Button>
            </>
          )}
        </div>
      </div>

      {/* Details card */}
      <div className="bg-white rounded-xl border shadow-sm p-5 space-y-4">
        <div className="flex items-center gap-2 text-sm text-gray-700">
          <ClockIcon className="h-5 w-5 text-gray-400" />
          <span>
            {event.is_all_day
              ? `All day · ${format(new Date(event.starts_at), 'MMMM d, yyyy')}`
              : `${format(new Date(event.starts_at), 'MMMM d, yyyy · h:mm a')} – ${format(new Date(event.ends_at), 'h:mm a')}`}
          </span>
        </div>

        {event.location && (
          <div className="flex items-center gap-2 text-sm text-gray-700">
            <MapPinIcon className="h-5 w-5 text-gray-400" />
            <span>{event.location}</span>
          </div>
        )}

        {event.description && (
          <p className="text-sm text-gray-600 leading-relaxed">{event.description}</p>
        )}
      </div>

      {/* My status */}
      <div className="bg-white rounded-xl border shadow-sm p-4">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            <span className="text-sm font-medium text-gray-700">Your status:</span>
            <StatusBadge status={myStatus} />
          </div>
          <div className="flex gap-1">
            {STATUS_OPTIONS.map((o) => (
              <Button
                key={o.value}
                variant={myStatus === o.value ? 'primary' : 'ghost'}
                size="sm"
                onClick={() => handleStatusChange(o.value)}
                loading={updateStatus.isPending}
              >
                {o.label}
              </Button>
            ))}
          </div>
        </div>
      </div>

      {/* Attendees */}
      {event.attendees && event.attendees.length > 0 && (
        <div className="bg-white rounded-xl border shadow-sm p-4">
          <div className="flex items-center justify-between mb-3">
            <h3 className="font-semibold text-gray-900 flex items-center gap-2">
              <UsersIcon className="h-5 w-5 text-gray-400" />
              Attendees ({event.attendees.length})
            </h3>
            {isOwner && (
              <Button variant="secondary" size="sm" onClick={() => setShowInviteModal(true)}>
                Invite
              </Button>
            )}
          </div>
          <div className="space-y-2">
            {event.attendees.map((a) => (
              <div key={a.user_id} className="flex items-center justify-between">
                <div className="flex items-center gap-2">
                  <div className="w-7 h-7 bg-blue-100 rounded-full flex items-center justify-center text-xs font-bold text-blue-700">
                    {a.user?.username?.[0]?.toUpperCase() ?? '?'}
                  </div>
                  <span className="text-sm text-gray-700">{a.user?.username ?? a.user_id}</span>
                  {a.is_organizer && (
                    <span className="text-xs bg-gray-100 text-gray-500 px-1.5 py-0.5 rounded-full">organizer</span>
                  )}
                </div>
                <StatusBadge status={a.status} />
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Invitations (organizer only) */}
      {isOwner && (
        <InvitationList eventId={event.id} onInvite={() => setShowInviteModal(true)} />
      )}

      <InviteUserModal
        open={showInviteModal}
        onClose={() => setShowInviteModal(false)}
        eventId={event.id}
      />
    </div>
  )
}
