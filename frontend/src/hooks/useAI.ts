import { useMutation } from '@tanstack/react-query'
import { generateDescription, parseEvent, suggestTimes } from '../api/ai'

export const useGenerateDescription = () =>
  useMutation({ mutationFn: generateDescription })

export const useParseEvent = () =>
  useMutation({ mutationFn: parseEvent })

export const useSuggestTimes = () =>
  useMutation({ mutationFn: suggestTimes })
