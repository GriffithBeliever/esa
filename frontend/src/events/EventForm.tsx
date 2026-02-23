import React, { useState, useEffect } from 'react'
import { format } from 'date-fns'
import type { Event, TimeSuggestion } from '../types'
import { Input, Textarea } from '../components/Input'
import { Button } from '../components/Button'
import { AIDescriptionPanel } from '../ai/AIDescriptionPanel'
import { ConflictWarning } from '../ai/ConflictWarning'
import { useConflicts } from '../hooks/useEvents'

interface EventFormProps {
  initial?: Partial<Event>
  onSubmit: (data: Partial<Event>) => Promise<void>
  loading?: boolean
  submitLabel?: string
}

const toDatetimeLocal = (iso: string) => {
  if (!iso) return ''
  const d = new Date(iso)
  return format(d, "yyyy-MM-dd'T'HH:mm")
}

const fromDatetimeLocal = (local: string) => {
  if (!local) return ''
  return new Date(local).toISOString()
}

export const EventForm: React.FC<EventFormProps> = ({
  initial = {},
  onSubmit,
  loading,
  submitLabel = 'Save event',
}) => {
  const [form, setForm] = useState({
    title: initial.title ?? '',
    description: initial.description ?? '',
    location: initial.location ?? '',
    starts_at: initial.starts_at ? toDatetimeLocal(initial.starts_at) : '',
    ends_at: initial.ends_at ? toDatetimeLocal(initial.ends_at) : '',
    is_all_day: initial.is_all_day ?? false,
    color: initial.color ?? '#3B82F6',
  })
  const [error, setError] = useState('')

  // Check conflicts whenever times change
  const conflictsEnabled =
    !!form.starts_at && !!form.ends_at && !form.is_all_day
  const { data: conflictData } = useConflicts(
    {
      starts_at: form.starts_at ? fromDatetimeLocal(form.starts_at) : '',
      ends_at: form.ends_at ? fromDatetimeLocal(form.ends_at) : '',
      exclude_id: initial.id,
    },
    conflictsEnabled
  )
  const conflicts = conflictData?.conflicts ?? []

  const handleSuggestion = (s: TimeSuggestion) => {
    setForm({
      ...form,
      starts_at: toDatetimeLocal(s.starts_at),
      ends_at: toDatetimeLocal(s.ends_at),
    })
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    if (!form.title.trim()) {
      setError('Title is required')
      return
    }
    if (!form.starts_at || !form.ends_at) {
      setError('Start and end times are required')
      return
    }
    try {
      await onSubmit({
        ...form,
        starts_at: fromDatetimeLocal(form.starts_at),
        ends_at: fromDatetimeLocal(form.ends_at),
      })
    } catch (err: unknown) {
      const msg = (err as { response?: { data?: { error?: string } } })?.response?.data?.error
      setError(msg || 'Failed to save event')
    }
  }

  // Color presets
  const colors = ['#3B82F6', '#10B981', '#F59E0B', '#EF4444', '#8B5CF6', '#EC4899', '#06B6D4', '#84CC16']

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      {error && (
        <div className="p-3 bg-red-50 border border-red-200 rounded-lg text-red-700 text-sm">
          {error}
        </div>
      )}

      <Input
        label="Title *"
        value={form.title}
        onChange={(e) => setForm({ ...form, title: e.target.value })}
        placeholder="Event title"
        required
        autoFocus
      />

      <div>
        <div className="flex items-center justify-between mb-1">
          <label className="block text-sm font-medium text-gray-700">Description</label>
          <AIDescriptionPanel
            title={form.title}
            location={form.location}
            onGenerated={(desc) => setForm({ ...form, description: desc })}
          />
        </div>
        <Textarea
          value={form.description}
          onChange={(e) => setForm({ ...form, description: e.target.value })}
          placeholder="What's this event about?"
          rows={3}
        />
      </div>

      <Input
        label="Location"
        value={form.location}
        onChange={(e) => setForm({ ...form, location: e.target.value })}
        placeholder="Where is it?"
      />

      <div className="flex items-center gap-2">
        <input
          type="checkbox"
          id="is_all_day"
          checked={form.is_all_day}
          onChange={(e) => setForm({ ...form, is_all_day: e.target.checked })}
          className="h-4 w-4 rounded border-gray-300 text-blue-600"
        />
        <label htmlFor="is_all_day" className="text-sm text-gray-700">All day event</label>
      </div>

      <div className="grid grid-cols-2 gap-4">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">
            {form.is_all_day ? 'Start date *' : 'Starts at *'}
          </label>
          <input
            type={form.is_all_day ? 'date' : 'datetime-local'}
            value={form.is_all_day ? form.starts_at.split('T')[0] : form.starts_at}
            onChange={(e) => {
              const val = form.is_all_day ? e.target.value + 'T00:00' : e.target.value
              setForm({ ...form, starts_at: val })
            }}
            required
            className="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">
            {form.is_all_day ? 'End date *' : 'Ends at *'}
          </label>
          <input
            type={form.is_all_day ? 'date' : 'datetime-local'}
            value={form.is_all_day ? form.ends_at.split('T')[0] : form.ends_at}
            onChange={(e) => {
              const val = form.is_all_day ? e.target.value + 'T23:59' : e.target.value
              setForm({ ...form, ends_at: val })
            }}
            required
            className="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>
      </div>

      {conflicts.length > 0 && (
        <ConflictWarning
          conflicts={conflicts}
          preferredStart={fromDatetimeLocal(form.starts_at)}
          onSelectSuggestion={handleSuggestion}
        />
      )}

      <div>
        <label className="block text-sm font-medium text-gray-700 mb-2">Color</label>
        <div className="flex gap-2 flex-wrap">
          {colors.map((c) => (
            <button
              key={c}
              type="button"
              onClick={() => setForm({ ...form, color: c })}
              className={`w-7 h-7 rounded-full border-2 transition-all ${
                form.color === c ? 'border-gray-800 scale-110' : 'border-transparent'
              }`}
              style={{ backgroundColor: c }}
            />
          ))}
        </div>
      </div>

      <Button type="submit" className="w-full" loading={loading}>
        {submitLabel}
      </Button>
    </form>
  )
}
