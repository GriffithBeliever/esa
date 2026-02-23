import client from './client'
import type { Event, EventFilter, PaginatedEvents, EventStatus } from '../types'

export const listEvents = (filter: EventFilter = {}) =>
  client.get<PaginatedEvents>('/events', { params: filter }).then((r) => r.data)

export const createEvent = (data: Partial<Event>) =>
  client.post<Event>('/events', data).then((r) => r.data)

export const getEvent = (id: string) =>
  client.get<Event>(`/events/${id}`).then((r) => r.data)

export const updateEvent = (id: string, data: Partial<Event>) =>
  client.put<Event>(`/events/${id}`, data).then((r) => r.data)

export const deleteEvent = (id: string) => client.delete(`/events/${id}`)

export const updateEventStatus = (id: string, status: EventStatus) =>
  client.patch<{ status: EventStatus }>(`/events/${id}/status`, { status }).then((r) => r.data)

export const checkConflicts = (params: { starts_at: string; ends_at: string; exclude_id?: string }) =>
  client.get<{ conflicts: Event[] }>('/events/conflicts', { params }).then((r) => r.data)
