import React, { useState } from 'react'
import { Modal } from '../components/Modal'
import { Input, Textarea } from '../components/Input'
import { Button } from '../components/Button'
import { useSendInvitation } from '../hooks/useInvitations'

interface InviteUserModalProps {
  open: boolean
  onClose: () => void
  eventId: string
}

export const InviteUserModal: React.FC<InviteUserModalProps> = ({ open, onClose, eventId }) => {
  const [email, setEmail] = useState('')
  const [message, setMessage] = useState('')
  const [success, setSuccess] = useState(false)
  const sendInvitation = useSendInvitation()

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    await sendInvitation.mutateAsync(
      { eventId, data: { invitee_email: email, message } },
      {
        onSuccess: () => {
          setSuccess(true)
          setEmail('')
          setMessage('')
          setTimeout(() => { setSuccess(false); onClose() }, 1500)
        },
      }
    )
  }

  return (
    <Modal open={open} onClose={onClose} title="Invite someone">
      {success ? (
        <div className="text-center py-6">
          <div className="text-green-500 text-4xl mb-2">✓</div>
          <p className="font-medium text-gray-900">Invitation sent!</p>
        </div>
      ) : (
        <form onSubmit={handleSubmit} className="space-y-4">
          <Input
            label="Email address *"
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            placeholder="friend@example.com"
            required
            autoFocus
          />
          <Textarea
            label="Personal message (optional)"
            value={message}
            onChange={(e) => setMessage(e.target.value)}
            placeholder="Hope to see you there!"
            rows={3}
          />
          {sendInvitation.isError && (
            <p className="text-sm text-red-600">Failed to send invitation</p>
          )}
          <div className="flex gap-2 justify-end">
            <Button type="button" variant="secondary" onClick={onClose}>
              Cancel
            </Button>
            <Button type="submit" loading={sendInvitation.isPending}>
              Send invitation
            </Button>
          </div>
        </form>
      )}
    </Modal>
  )
}
