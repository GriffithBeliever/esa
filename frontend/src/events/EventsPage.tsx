import React, { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useEvents, useDeleteEvent } from '../hooks/useEvents'
import { EventCard } from './EventCard'
import { EventSearch } from './EventSearch'
import { CalendarView } from './CalendarView'
import { NaturalLanguageInput } from '../ai/NaturalLanguageInput'
import { Spinner } from '../components/Spinner'
import { Button } from '../components/Button'
import type { EventFilter, ParsedEvent } from '../types'
import { CalendarDaysIcon, ListBulletIcon } from '@heroicons/react/24/outline'
import { Modal } from '../components/Modal'
import { EventForm } from './EventForm'
import { useCreateEvent } from '../hooks/useEvents'

type ViewMode = 'list' | 'calendar'

export const EventsPage: React.FC = () => {
  const navigate = useNavigate()
  const [filter, setFilter] = useState<EventFilter>({ page: 1, limit: 20 })
  const [view, setView] = useState<ViewMode>('list')
  const [showQuickAdd, setShowQuickAdd] = useState(false)
  const [parsedEvent, setParsedEvent] = useState<Partial<{
    title: string; location: string; starts_at: string; ends_at: string
  }> | null>(null)

  const { data, isLoading } = useEvents(filter)
  const deleteEvent = useDeleteEvent()
  const createEvent = useCreateEvent()

  const handleParsed = (parsed: ParsedEvent) => {
    setParsedEvent({
      title: parsed.title,
      location: parsed.location,
      starts_at: parsed.starts_at,
      ends_at: parsed.ends_at,
    })
    setShowQuickAdd(true)
  }

  const handleQuickCreate = async (formData: Record<string, unknown>) => {
    await createEvent.mutateAsync(formData)
    setShowQuickAdd(false)
    setParsedEvent(null)
  }

  const handleDelete = async (id: string) => {
    if (!confirm('Delete this event?')) return
    await deleteEvent.mutateAsync(id)
  }

  const events = data?.events ?? []
  const total = data?.total ?? 0

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold text-gray-900">Events</h1>
        <div className="flex items-center gap-2">
          <button
            onClick={() => setView('list')}
            className={`p-2 rounded-lg ${view === 'list' ? 'bg-blue-100 text-blue-700' : 'text-gray-500 hover:bg-gray-100'}`}
            title="List view"
          >
            <ListBulletIcon className="h-5 w-5" />
          </button>
          <button
            onClick={() => setView('calendar')}
            className={`p-2 rounded-lg ${view === 'calendar' ? 'bg-blue-100 text-blue-700' : 'text-gray-500 hover:bg-gray-100'}`}
            title="Calendar view"
          >
            <CalendarDaysIcon className="h-5 w-5" />
          </button>
          <Button onClick={() => navigate('/events/new')}>New Event</Button>
        </div>
      </div>

      {/* Natural Language Quick-Add */}
      <NaturalLanguageInput onParsed={handleParsed} />

      {/* Filters (list view only) */}
      {view === 'list' && (
        <EventSearch filter={filter} onChange={setFilter} />
      )}

      {/* Content */}
      {isLoading ? (
        <div className="flex justify-center py-12"><Spinner /></div>
      ) : view === 'calendar' ? (
        <CalendarView events={events} />
      ) : events.length === 0 ? (
        <div className="text-center py-16 text-gray-500">
          <CalendarDaysIcon className="h-12 w-12 mx-auto mb-3 text-gray-300" />
          <p className="font-medium">No events found</p>
          <p className="text-sm mt-1">Create your first event or try a different search</p>
          <Button className="mt-4" onClick={() => navigate('/events/new')}>
            Create event
          </Button>
        </div>
      ) : (
        <>
          <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4">
            {events.map((e) => (
              <EventCard key={e.id} event={e} onDelete={handleDelete} />
            ))}
          </div>
          {/* Pagination */}
          {total > (filter.limit ?? 20) && (
            <div className="flex items-center justify-between py-4">
              <p className="text-sm text-gray-500">
                Showing {((filter.page ?? 1) - 1) * (filter.limit ?? 20) + 1}–
                {Math.min((filter.page ?? 1) * (filter.limit ?? 20), total)} of {total}
              </p>
              <div className="flex gap-2">
                <Button
                  variant="secondary"
                  size="sm"
                  disabled={(filter.page ?? 1) <= 1}
                  onClick={() => setFilter({ ...filter, page: (filter.page ?? 1) - 1 })}
                >
                  Previous
                </Button>
                <Button
                  variant="secondary"
                  size="sm"
                  disabled={(filter.page ?? 1) * (filter.limit ?? 20) >= total}
                  onClick={() => setFilter({ ...filter, page: (filter.page ?? 1) + 1 })}
                >
                  Next
                </Button>
              </div>
            </div>
          )}
        </>
      )}

      {/* Quick-add modal (pre-filled from NL parse) */}
      <Modal
        open={showQuickAdd}
        onClose={() => { setShowQuickAdd(false); setParsedEvent(null) }}
        title="Create event"
        size="lg"
      >
        <EventForm
          initial={parsedEvent ?? {}}
          onSubmit={handleQuickCreate}
          loading={createEvent.isPending}
          submitLabel="Create event"
        />
      </Modal>
    </div>
  )
}
