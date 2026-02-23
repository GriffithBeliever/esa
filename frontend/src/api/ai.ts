import client from './client'
import type { ParsedEvent, TimeSuggestion, Event } from '../types'

export const generateDescription = (data: { title: string; location: string }) =>
  client.post<{ description: string }>('/ai/generate-description', data).then((r) => r.data)

export const parseEvent = (data: { text: string; today?: string }) =>
  client.post<ParsedEvent>('/ai/parse-event', data).then((r) => r.data)

export const suggestTimes = (data: { preferred_time: string; conflicts: Event[] }) =>
  client.post<{ suggestions: TimeSuggestion[] }>('/ai/suggest-times', data).then((r) => r.data)
