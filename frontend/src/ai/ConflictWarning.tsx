import React from 'react'
import { format } from 'date-fns'
import { ExclamationTriangleIcon, LightBulbIcon } from '@heroicons/react/24/outline'
import type { Event, TimeSuggestion } from '../types'
import { useSuggestTimes } from '../hooks/useAI'
import { Button } from '../components/Button'

interface ConflictWarningProps {
  conflicts: Event[]
  preferredStart: string
  onSelectSuggestion: (suggestion: TimeSuggestion) => void
}

export const ConflictWarning: React.FC<ConflictWarningProps> = ({
  conflicts,
  preferredStart,
  onSelectSuggestion,
}) => {
  const { mutate, data, isPending } = useSuggestTimes()

  if (conflicts.length === 0) return null

  const getSuggestions = () => {
    mutate({
      preferred_time: preferredStart,
      conflicts,
    })
  }

  return (
    <div className="border border-amber-200 bg-amber-50 rounded-lg p-4 space-y-3">
      <div className="flex items-start gap-2">
        <ExclamationTriangleIcon className="h-5 w-5 text-amber-600 flex-shrink-0 mt-0.5" />
        <div>
          <p className="text-sm font-medium text-amber-800">
            Scheduling conflict detected
          </p>
          <ul className="mt-1 space-y-0.5">
            {conflicts.map((c) => (
              <li key={c.id} className="text-xs text-amber-700">
                {c.title} · {format(new Date(c.starts_at), 'h:mm a')} – {format(new Date(c.ends_at), 'h:mm a')}
              </li>
            ))}
          </ul>
        </div>
      </div>

      <Button
        type="button"
        variant="ghost"
        size="sm"
        onClick={getSuggestions}
        loading={isPending}
        className="text-purple-600 hover:bg-purple-50"
      >
        <LightBulbIcon className="h-4 w-4 mr-1" />
        Suggest alternative times
      </Button>

      {data?.suggestions && data.suggestions.length > 0 && (
        <div className="space-y-2">
          <p className="text-xs font-medium text-gray-600">AI suggestions:</p>
          {data.suggestions.map((s, i) => (
            <div key={i} className="bg-white rounded-lg border p-3">
              <div className="flex items-center justify-between mb-1">
                <span className="text-sm font-medium text-gray-900">
                  {format(new Date(s.starts_at), 'MMM d, h:mm a')} – {format(new Date(s.ends_at), 'h:mm a')}
                </span>
                <Button
                  type="button"
                  variant="secondary"
                  size="sm"
                  onClick={() => onSelectSuggestion(s)}
                >
                  Use this
                </Button>
              </div>
              <p className="text-xs text-gray-500">{s.reasoning}</p>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
