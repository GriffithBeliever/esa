import React from 'react'
import { useNavigate } from 'react-router-dom'
import { useCreateEvent } from '../hooks/useEvents'
import { EventForm } from './EventForm'
import type { Event } from '../types'

export const CreateEventPage: React.FC = () => {
  const navigate = useNavigate()
  const createEvent = useCreateEvent()

  const handleSubmit = async (data: Partial<Event>) => {
    const event = await createEvent.mutateAsync(data)
    navigate(`/events/${event.id}`)
  }

  return (
    <div className="max-w-2xl mx-auto">
      <h1 className="text-2xl font-bold text-gray-900 mb-6">Create new event</h1>
      <div className="bg-white rounded-xl border shadow-sm p-6">
        <EventForm
          onSubmit={handleSubmit}
          loading={createEvent.isPending}
          submitLabel="Create event"
        />
      </div>
    </div>
  )
}
