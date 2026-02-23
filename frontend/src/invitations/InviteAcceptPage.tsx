import React, { useState } from 'react'
import { useParams, Link } from 'react-router-dom'
import { respondToInvitation } from '../api/invitations'
import { Button } from '../components/Button'
import { CalendarDaysIcon } from '@heroicons/react/24/outline'

export const InviteAcceptPage: React.FC = () => {
  const { token } = useParams<{ token: string }>()
  const [status, setStatus] = useState<'idle' | 'loading' | 'accepted' | 'declined' | 'error'>('idle')
  const [errMsg, setErrMsg] = useState('')

  const respond = async (accept: boolean) => {
    if (!token) return
    setStatus('loading')
    try {
      await respondToInvitation(token, accept)
      setStatus(accept ? 'accepted' : 'declined')
    } catch (err: unknown) {
      const msg = (err as { response?: { data?: { error?: string } } })?.response?.data?.error
      setErrMsg(msg || 'Something went wrong')
      setStatus('error')
    }
  }

  return (
    <div className="min-h-screen bg-gray-50 flex items-center justify-center p-4">
      <div className="bg-white rounded-2xl shadow-sm border w-full max-w-md p-8 text-center">
        <CalendarDaysIcon className="h-12 w-12 mx-auto text-blue-600 mb-4" />

        {status === 'accepted' && (
          <>
            <h1 className="text-2xl font-bold text-green-700 mb-2">You're in! 🎉</h1>
            <p className="text-gray-500 mb-4">Your response has been recorded.</p>
            <Link to="/events" className="text-blue-600 hover:underline text-sm">
              View your events →
            </Link>
          </>
        )}

        {status === 'declined' && (
          <>
            <h1 className="text-2xl font-bold text-gray-700 mb-2">Invitation declined</h1>
            <p className="text-gray-500 mb-4">Your response has been recorded.</p>
            <Link to="/" className="text-blue-600 hover:underline text-sm">
              Go home →
            </Link>
          </>
        )}

        {status === 'error' && (
          <>
            <h1 className="text-2xl font-bold text-red-700 mb-2">Something went wrong</h1>
            <p className="text-gray-500 mb-4">{errMsg}</p>
          </>
        )}

        {(status === 'idle' || status === 'loading') && (
          <>
            <h1 className="text-2xl font-bold text-gray-900 mb-2">You've been invited!</h1>
            <p className="text-gray-500 mb-6">Would you like to attend this event?</p>
            <div className="flex gap-3 justify-center">
              <Button
                onClick={() => respond(true)}
                loading={status === 'loading'}
                size="lg"
              >
                Accept
              </Button>
              <Button
                variant="secondary"
                onClick={() => respond(false)}
                loading={status === 'loading'}
                size="lg"
              >
                Decline
              </Button>
            </div>
          </>
        )}
      </div>
    </div>
  )
}
