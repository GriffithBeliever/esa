import React, { useState } from 'react'
import { SparklesIcon, ArrowRightIcon } from '@heroicons/react/24/outline'
import { useParseEvent } from '../hooks/useAI'
import type { ParsedEvent } from '../types'
import { format } from 'date-fns'

interface NaturalLanguageInputProps {
  onParsed: (event: ParsedEvent) => void
}

export const NaturalLanguageInput: React.FC<NaturalLanguageInputProps> = ({ onParsed }) => {
  const [text, setText] = useState('')
  const { mutate, isPending, error } = useParseEvent()

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (!text.trim()) return
    mutate(
      { text, today: format(new Date(), 'yyyy-MM-dd') },
      {
        onSuccess: (data) => {
          onParsed(data)
          setText('')
        },
      }
    )
  }

  return (
    <form onSubmit={handleSubmit} className="flex gap-2 items-center bg-white border rounded-xl px-4 py-3 shadow-sm">
      <SparklesIcon className="h-5 w-5 text-purple-500 flex-shrink-0" />
      <input
        type="text"
        value={text}
        onChange={(e) => setText(e.target.value)}
        placeholder="Quick add: &quot;team lunch next Tuesday at noon&quot;"
        className="flex-1 text-sm bg-transparent focus:outline-none placeholder-gray-400"
      />
      <button
        type="submit"
        disabled={!text.trim() || isPending}
        className="flex items-center gap-1 px-3 py-1.5 bg-purple-600 text-white rounded-lg text-sm font-medium hover:bg-purple-700 disabled:opacity-50 disabled:cursor-not-allowed"
      >
        {isPending ? (
          <svg className="animate-spin h-4 w-4" fill="none" viewBox="0 0 24 24">
            <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
            <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
          </svg>
        ) : (
          <>
            Parse <ArrowRightIcon className="h-3 w-3" />
          </>
        )}
      </button>
      {error && <span className="text-xs text-red-500">Parse failed</span>}
    </form>
  )
}
