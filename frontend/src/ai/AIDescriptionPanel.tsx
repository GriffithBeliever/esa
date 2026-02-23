import React, { useState } from 'react'
import { SparklesIcon } from '@heroicons/react/24/outline'
import { useGenerateDescription } from '../hooks/useAI'
import { Button } from '../components/Button'

interface AIDescriptionPanelProps {
  title: string
  location: string
  onGenerated: (desc: string) => void
}

export const AIDescriptionPanel: React.FC<AIDescriptionPanelProps> = ({ title, location, onGenerated }) => {
  const { mutate, isPending, error } = useGenerateDescription()
  const [errMsg, setErrMsg] = useState('')

  const generate = () => {
    if (!title) {
      setErrMsg('Enter a title first')
      return
    }
    setErrMsg('')
    mutate(
      { title, location },
      {
        onSuccess: (data) => onGenerated(data.description),
        onError: () => setErrMsg('AI generation failed'),
      }
    )
  }

  return (
    <div className="flex items-center gap-2">
      <Button
        type="button"
        variant="ghost"
        size="sm"
        onClick={generate}
        loading={isPending}
        className="text-purple-600 hover:bg-purple-50"
      >
        <SparklesIcon className="h-4 w-4 mr-1" />
        Generate with AI
      </Button>
      {(error || errMsg) && (
        <span className="text-xs text-red-500">{errMsg || 'Failed'}</span>
      )}
    </div>
  )
}
