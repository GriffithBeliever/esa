import React from 'react'
import { useIncomingInvitations } from '../hooks/useInvitations'
import { InvitationCard } from './InvitationCard'
import { Spinner } from '../components/Spinner'
import { BellIcon } from '@heroicons/react/24/outline'

export const InvitationsPage: React.FC = () => {
  const { data, isLoading } = useIncomingInvitations()
  const invitations = data?.invitations ?? []

  return (
    <div className="max-w-2xl mx-auto space-y-4">
      <h1 className="text-2xl font-bold text-gray-900">Invitations</h1>

      {isLoading ? (
        <div className="flex justify-center py-12"><Spinner /></div>
      ) : invitations.length === 0 ? (
        <div className="text-center py-16 text-gray-500">
          <BellIcon className="h-12 w-12 mx-auto mb-3 text-gray-300" />
          <p className="font-medium">No pending invitations</p>
          <p className="text-sm mt-1">When someone invites you to an event, it will appear here</p>
        </div>
      ) : (
        <div className="space-y-3">
          {invitations.map((inv) => (
            <InvitationCard key={inv.id} invitation={inv} />
          ))}
        </div>
      )}
    </div>
  )
}
