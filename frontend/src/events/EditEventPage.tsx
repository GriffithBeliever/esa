import React from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { useQuery } from '@tanstack/react-query'
import { getEvent } from '../api/events'
import { useUpdateEvent } from '../hooks/useEvents'
import { EventForm } from './EventForm'
import { Spinner } from '../components/Spinner'
import type { Event } from '../types'

export const EditEventPage: React.FC = () => {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const updateEvent = useUpdateEvent()

  const { data: event, isLoading } = useQuery({
    queryKey: ['event', id],
    queryFn: () => getEvent(id!),
    enabled: !!id,
  })

  if (isLoading) return <div className="flex justify-center p-12"><Spinner /></div>
  if (!event) return <div className="text-center p-12 text-gray-500">Event not found</div>

  const handleSubmit = async (data: Partial<Event>) => {
    await updateEvent.mutateAsync({ id: event.id, data })
    navigate(`/events/${event.id}`)
  }

  return (
    <div className="max-w-2xl mx-auto">
      <h1 className="text-2xl font-bold text-gray-900 mb-6">Edit event</h1>
      <div className="bg-white rounded-xl border shadow-sm p-6">
        <EventForm
          initial={event}
          onSubmit={handleSubmit}
          loading={updateEvent.isPending}
          submitLabel="Save changes"
        />
      </div>
    </div>
  )
}
