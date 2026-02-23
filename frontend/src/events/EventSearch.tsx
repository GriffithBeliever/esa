import React, { useState } from 'react'
import { MagnifyingGlassIcon, FunnelIcon } from '@heroicons/react/24/outline'
import type { EventFilter, EventStatus } from '../types'
import { Input } from '../components/Input'

interface EventSearchProps {
  filter: EventFilter
  onChange: (filter: EventFilter) => void
}

const STATUS_OPTIONS: { value: EventStatus | ''; label: string }[] = [
  { value: '', label: 'All statuses' },
  { value: 'upcoming', label: 'Upcoming' },
  { value: 'attending', label: 'Attending' },
  { value: 'maybe', label: 'Maybe' },
  { value: 'declined', label: 'Declined' },
]

export const EventSearch: React.FC<EventSearchProps> = ({ filter, onChange }) => {
  const [expanded, setExpanded] = useState(false)

  return (
    <div className="bg-white border rounded-xl p-4 space-y-3">
      <div className="flex gap-2">
        <div className="relative flex-1">
          <MagnifyingGlassIcon className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-gray-400" />
          <input
            type="text"
            value={filter.q ?? ''}
            onChange={(e) => onChange({ ...filter, q: e.target.value, page: 1 })}
            placeholder="Search events..."
            className="w-full pl-9 pr-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>
        <button
          onClick={() => setExpanded(!expanded)}
          className={`flex items-center gap-1.5 px-3 py-2 rounded-lg border text-sm ${
            expanded ? 'bg-blue-50 border-blue-200 text-blue-700' : 'border-gray-300 text-gray-600'
          }`}
        >
          <FunnelIcon className="h-4 w-4" />
          Filters
        </button>
      </div>

      {expanded && (
        <div className="grid grid-cols-2 gap-3 pt-2 border-t">
          <div>
            <label className="block text-xs font-medium text-gray-600 mb-1">Status</label>
            <select
              value={filter.status ?? ''}
              onChange={(e) => onChange({ ...filter, status: (e.target.value as EventStatus) || undefined, page: 1 })}
              className="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              {STATUS_OPTIONS.map((o) => (
                <option key={o.value} value={o.value}>{o.label}</option>
              ))}
            </select>
          </div>

          <Input
            label="Location"
            value={filter.location ?? ''}
            onChange={(e) => onChange({ ...filter, location: e.target.value, page: 1 })}
            placeholder="Filter by location"
          />

          <div>
            <label className="block text-xs font-medium text-gray-600 mb-1">From</label>
            <input
              type="date"
              value={filter.from ? filter.from.split('T')[0] : ''}
              onChange={(e) => onChange({ ...filter, from: e.target.value ? new Date(e.target.value).toISOString() : undefined, page: 1 })}
              className="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>

          <div>
            <label className="block text-xs font-medium text-gray-600 mb-1">To</label>
            <input
              type="date"
              value={filter.to ? filter.to.split('T')[0] : ''}
              onChange={(e) => onChange({ ...filter, to: e.target.value ? new Date(e.target.value).toISOString() : undefined, page: 1 })}
              className="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>
        </div>
      )}
    </div>
  )
}
