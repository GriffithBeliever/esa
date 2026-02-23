import React from 'react'
import { format } from 'date-fns'
import { useEventInvitations, useDeleteInvitation } from '../hooks/useInvitations'
import { Button } from '../components/Button'
import { Spinner } from '../components/Spinner'
import { EnvelopeIcon, TrashIcon } from '@heroicons/react/24/outline'

interface InvitationListProps {
  eventId: string
  onInvite: () => void
}

const statusColors: Record<string, string> = {
  pending: 'text-amber-600',
  accepted: 'text-green-600',
  declined: 'text-gray-500',
  expired: 'text-red-500',
}

export const InvitationList: React.FC<InvitationListProps> = ({ eventId, onInvite }) => {
  const { data, isLoading } = useEventInvitations(eventId)
  const deleteInv = useDeleteInvitation()

  const invitations = data?.invitations ?? []

  return (
    <div className="bg-white rounded-xl border shadow-sm p-4">
      <div className="flex items-center justify-between mb-4">
        <h3 className="font-semibold text-gray-900 flex items-center gap-2">
          <EnvelopeIcon className="h-5 w-5 text-gray-400" />
          Invitations ({invitations.length})
        </h3>
        <Button variant="secondary" size="sm" onClick={onInvite}>
          Send invitation
        </Button>
      </div>

      {isLoading ? (
        <Spinner size="sm" />
      ) : invitations.length === 0 ? (
        <p className="text-sm text-gray-500 text-center py-4">No invitations sent yet</p>
      ) : (
        <div className="space-y-2">
          {invitations.map((inv) => (
            <div key={inv.id} className="flex items-center justify-between py-2 border-b last:border-0">
              <div>
                <p className="text-sm font-medium text-gray-900">{inv.invitee_email}</p>
                <p className="text-xs text-gray-500">
                  Sent {format(new Date(inv.created_at), 'MMM d')} ·{' '}
                  <span className={`font-medium ${statusColors[inv.status]}`}>{inv.status}</span>
                </p>
              </div>
              <button
                onClick={() => deleteInv.mutate({ eventId, invId: inv.id })}
                className="p-1.5 rounded hover:bg-red-50 text-gray-400 hover:text-red-500"
                title="Cancel invitation"
              >
                <TrashIcon className="h-4 w-4" />
              </button>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
