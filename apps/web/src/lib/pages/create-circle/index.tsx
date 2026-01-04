import { useState } from 'react'
import { useNavigate } from '@tanstack/react-router'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { toast } from 'react-hot-toast'

import { api_url } from '@/lib/constants'

const CreateGroup = () => {
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const [name, setName] = useState('')
  const [contributionAmount, setContributionAmount] = useState('')
  const [frequency, setFrequency] = useState('Weekly')
  const [isPublic, setIsPublic] = useState(false)

  const createMutation = useMutation({
    mutationFn: async (data: { name: string; description?: string }) => {
      const res = await fetch(`${api_url}/groups`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify(data),
      })
      if (!res.ok) {
        const error = await res
          .json()
          .catch(() => ({ message: 'Failed to create group' }))
        throw new Error(error.message || 'Failed to create group')
      }
      return res.json()
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['/groups'] })
      toast.success('Group created successfully')
      navigate({ to: '/app' })
    },
    onError: (error: Error) => {
      toast.error(error.message || 'Failed to create group')
    },
  })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    createMutation.mutate({
      name,
    })
  }

  return (
    <div className="max-w-2xl mx-auto pt-12">
      <div className="bg-white dark:bg-[#241f2e] border border-[#eeeeee] dark:border-[#3d3551] rounded-xl shadow-sm p-6 md:p-8">
        <h2 className="text-xl font-semibold tracking-tight text-[#333] dark:text-[#f5f3ff] mb-1">
          Create a Circle
        </h2>
        <p className="text-sm text-[#666] dark:text-[#c4b5fd] mb-8">
          Set up a private ROSCA group. Invite friends after creation.
        </p>

        <form onSubmit={handleSubmit} className="space-y-6">
          <div>
            <label className="block text-xs font-medium text-[#666] dark:text-[#c4b5fd] mb-1.5 uppercase tracking-wide">
              Group Name
            </label>
            <input
              type="text"
              required
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="e.g. Summer Vacation Fund"
              className="w-full px-3 py-2.5 bg-white dark:bg-[#241f2e] border border-[#eeeeee] dark:border-[#3d3551] rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-[#667eea]/20 dark:focus:ring-[#c4b5fd]/20 focus:border-[#667eea] dark:focus:border-[#c4b5fd] transition-all placeholder:text-[#999] dark:placeholder:text-[#a78bfa] text-[#333] dark:text-[#f5f3ff]"
            />
          </div>

          <div>
            <label className="block text-xs font-medium text-[#666] dark:text-[#c4b5fd] mb-1.5 uppercase tracking-wide">
              Contribution Amount
            </label>
            <div className="relative">
              <span className="absolute left-3 top-2.5 text-[#999] dark:text-[#a78bfa] text-sm">
                $
              </span>
              <input
                type="number"
                required
                value={contributionAmount}
                onChange={(e) => setContributionAmount(e.target.value)}
                placeholder="100"
                className="w-full pl-7 pr-3 py-2.5 bg-white dark:bg-[#241f2e] border border-[#eeeeee] dark:border-[#3d3551] rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-[#667eea]/20 dark:focus:ring-[#c4b5fd]/20 focus:border-[#667eea] dark:focus:border-[#c4b5fd] transition-all placeholder:text-[#999] dark:placeholder:text-[#a78bfa] text-[#333] dark:text-[#f5f3ff]"
              />
            </div>
          </div>

          <div>
            <label className="block text-xs font-medium text-[#666] dark:text-[#c4b5fd] mb-1.5 uppercase tracking-wide">
              Payout Frequency
            </label>
            <div className="grid grid-cols-3 gap-3">
              <label className="cursor-pointer">
                <input
                  type="radio"
                  name="freq"
                  value="Weekly"
                  checked={frequency === 'Weekly'}
                  onChange={(e) => setFrequency(e.target.value)}
                  className="peer sr-only"
                />
                <div
                  className={`text-center py-2.5 border rounded-lg text-sm transition-all hover:bg-[#f5f5f5] dark:hover:bg-[#3d3551] ${
                    frequency === 'Weekly'
                      ? 'bg-gradient-primary text-white border-[#667eea] dark:border-[#c4b5fd]'
                      : 'border-[#eeeeee] dark:border-[#3d3551] text-[#333] dark:text-[#c4b5fd]'
                  }`}
                >
                  Weekly
                </div>
              </label>
              <label className="cursor-pointer">
                <input
                  type="radio"
                  name="freq"
                  value="Bi-Weekly"
                  checked={frequency === 'Bi-Weekly'}
                  onChange={(e) => setFrequency(e.target.value)}
                  className="peer sr-only"
                />
                <div
                  className={`text-center py-2.5 border rounded-lg text-sm transition-all hover:bg-[#f5f5f5] dark:hover:bg-[#3d3551] ${
                    frequency === 'Bi-Weekly'
                      ? 'bg-gradient-primary text-white border-[#667eea] dark:border-[#c4b5fd]'
                      : 'border-[#eeeeee] dark:border-[#3d3551] text-[#333] dark:text-[#c4b5fd]'
                  }`}
                >
                  Bi-Weekly
                </div>
              </label>
              <label className="cursor-pointer">
                <input
                  type="radio"
                  name="freq"
                  value="Monthly"
                  checked={frequency === 'Monthly'}
                  onChange={(e) => setFrequency(e.target.value)}
                  className="peer sr-only"
                />
                <div
                  className={`text-center py-2.5 border rounded-lg text-sm transition-all hover:bg-[#f5f5f5] dark:hover:bg-[#3d3551] ${
                    frequency === 'Monthly'
                      ? 'bg-gradient-primary text-white border-[#667eea] dark:border-[#c4b5fd]'
                      : 'border-[#eeeeee] dark:border-[#3d3551] text-[#333] dark:text-[#c4b5fd]'
                  }`}
                >
                  Monthly
                </div>
              </label>
            </div>
          </div>

          <div className="pt-4 border-t border-[#eeeeee] dark:border-[#3d3551] flex items-center justify-between">
            <div className="flex items-center gap-3">
              <div className="relative inline-block w-10 h-6 align-middle select-none transition duration-200 ease-in">
                <input
                  type="checkbox"
                  id="toggle"
                  checked={isPublic}
                  onChange={(e) => setIsPublic(e.target.checked)}
                  className="toggle-checkbox absolute block w-4 h-4 rounded-full bg-white dark:bg-[#f5f3ff] border-2 border-[#eeeeee] dark:border-[#3d3551] appearance-none cursor-pointer top-1 left-1 transition-all duration-300"
                />
                <label
                  htmlFor="toggle"
                  className="toggle-label block overflow-hidden h-6 rounded-full cursor-pointer transition-colors duration-300"
                />
              </div>
              <label
                htmlFor="toggle"
                className="text-sm text-[#666] dark:text-[#c4b5fd] cursor-pointer select-none"
              >
                Publicly Visible (Discovery)
              </label>
            </div>
          </div>

          <div className="pt-2">
            <button
              type="submit"
              disabled={createMutation.isPending || !name.trim()}
              className="w-full bg-gradient-primary hover:opacity-90 disabled:opacity-50 disabled:cursor-not-allowed text-white font-medium text-sm py-3 rounded-lg transition-all shadow-md shadow-[#667eea]/10 flex items-center justify-center gap-2"
            >
              {createMutation.isPending ? (
                <>
                  <div className="loader ease-linear rounded-full border-2 border-t-2 border-white h-4 w-4" />
                  Creating Group...
                </>
              ) : (
                'Create Group'
              )}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}

export default CreateGroup
