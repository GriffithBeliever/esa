import client from './client'
import type { Invitation } from '../types'

export const sendInvitation = (eventId: string, data: { invitee_email: string; message?: string }) =>
  client.post<Invitation>(`/events/${eventId}/invitations`, data).then((r) => r.data)

export const listEventInvitations = (eventId: string) =>
  client.get<{ invitations: Invitation[] }>(`/events/${eventId}/invitations`).then((r) => r.data)

export const deleteInvitation = (eventId: string, invId: string) =>
  client.delete(`/events/${eventId}/invitations/${invId}`)

export const listIncomingInvitations = () =>
  client.get<{ invitations: Invitation[] }>('/invitations/incoming').then((r) => r.data)

export const respondToInvitation = (token: string, accept: boolean) =>
  client.post<{ status: string }>(`/invitations/${token}/respond`, { accept }).then((r) => r.data)
