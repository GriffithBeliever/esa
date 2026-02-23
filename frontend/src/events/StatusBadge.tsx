import type { EventStatus } from '../types'

const styles: Record<EventStatus, string> = {
  upcoming: 'bg-blue-100 text-blue-700',
  attending: 'bg-green-100 text-green-700',
  maybe: 'bg-yellow-100 text-yellow-700',
  declined: 'bg-gray-100 text-gray-500',
}

const labels: Record<EventStatus, string> = {
  upcoming: 'Upcoming',
  attending: 'Attending',
  maybe: 'Maybe',
  declined: 'Declined',
}

export const StatusBadge = ({ status }: { status: EventStatus }) => (
  <span className={`text-xs font-medium px-2 py-0.5 rounded-full ${styles[status]}`}>
    {labels[status]}
  </span>
)
