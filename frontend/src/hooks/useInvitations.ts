import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import {
  sendInvitation,
  listEventInvitations,
  listIncomingInvitations,
  respondToInvitation,
  deleteInvitation,
} from '../api/invitations'

export const useEventInvitations = (eventId: string) =>
  useQuery({
    queryKey: ['invitations', 'event', eventId],
    queryFn: () => listEventInvitations(eventId),
    enabled: !!eventId,
  })

export const useIncomingInvitations = () =>
  useQuery({
    queryKey: ['invitations', 'incoming'],
    queryFn: listIncomingInvitations,
  })

export const useSendInvitation = () => {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ eventId, data }: { eventId: string; data: Parameters<typeof sendInvitation>[1] }) =>
      sendInvitation(eventId, data),
    onSuccess: (_data, { eventId }) =>
      qc.invalidateQueries({ queryKey: ['invitations', 'event', eventId] }),
  })
}

export const useRespondInvitation = () => {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ token, accept }: { token: string; accept: boolean }) =>
      respondToInvitation(token, accept),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['invitations', 'incoming'] })
      qc.invalidateQueries({ queryKey: ['events'] })
    },
  })
}

export const useDeleteInvitation = () => {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ eventId, invId }: { eventId: string; invId: string }) =>
      deleteInvitation(eventId, invId),
    onSuccess: (_data, { eventId }) =>
      qc.invalidateQueries({ queryKey: ['invitations', 'event', eventId] }),
  })
}
