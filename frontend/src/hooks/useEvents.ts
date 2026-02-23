import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { listEvents, createEvent, updateEvent, deleteEvent, updateEventStatus, checkConflicts } from '../api/events'
import type { EventFilter, EventStatus } from '../types'

export const useEvents = (filter: EventFilter = {}) =>
  useQuery({
    queryKey: ['events', filter],
    queryFn: () => listEvents(filter),
  })

export const useCreateEvent = () => {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: createEvent,
    onSuccess: () => qc.invalidateQueries({ queryKey: ['events'] }),
  })
}

export const useUpdateEvent = () => {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: Parameters<typeof updateEvent>[1] }) =>
      updateEvent(id, data),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['events'] }),
  })
}

export const useDeleteEvent = () => {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: deleteEvent,
    onSuccess: () => qc.invalidateQueries({ queryKey: ['events'] }),
  })
}

export const useUpdateEventStatus = () => {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, status }: { id: string; status: EventStatus }) => updateEventStatus(id, status),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['events'] }),
  })
}

export const useConflicts = (params: { starts_at: string; ends_at: string; exclude_id?: string }, enabled = true) =>
  useQuery({
    queryKey: ['conflicts', params],
    queryFn: () => checkConflicts(params),
    enabled: enabled && !!params.starts_at && !!params.ends_at,
  })
