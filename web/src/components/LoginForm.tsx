import React, { useState } from 'react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { useSendLoginLink } from '@/hooks/useAuth'

export function LoginForm() {
  const [email, setEmail] = useState('')
  const sendLoginLinkMutation = useSendLoginLink()

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    sendLoginLinkMutation.mutate(email, {
      onSuccess: () => {
        setEmail('')
      }
    })
  }

  return (
    <Card className="w-full max-w-md mx-auto">
      <CardHeader>
        <CardTitle>Sign In</CardTitle>
        <CardDescription>
          Enter your email address and we'll send you a secure login link
        </CardDescription>
      </CardHeader>
      <CardContent>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="email">Email Address</Label>
            <Input
              id="email"
              type="email"
              placeholder="your@email.com"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
              disabled={sendLoginLinkMutation.isPending}
            />
          </div>

          {sendLoginLinkMutation.isSuccess && (
            <div className="p-3 text-sm text-green-800 bg-green-100 border border-green-200 rounded">
              Login link sent! Check your email and click the link to sign in.
            </div>
          )}

          {sendLoginLinkMutation.isError && (
            <div className="p-3 text-sm text-red-800 bg-red-100 border border-red-200 rounded">
              {sendLoginLinkMutation.error.message}
            </div>
          )}

          <Button 
            type="submit" 
            className="w-full" 
            disabled={sendLoginLinkMutation.isPending || !email}
          >
            {sendLoginLinkMutation.isPending ? 'Sending...' : 'Send Login Link'}
          </Button>
        </form>

        <div className="mt-6 text-sm text-gray-600">
          <p className="mb-2"><strong>How it works:</strong></p>
          <ol className="list-decimal list-inside space-y-1 text-xs">
            <li>Enter your email address above</li>
            <li>We'll send you a secure login link</li>
            <li>Click the link in your email to sign in</li>
            <li>No password required!</li>
          </ol>
        </div>
      </CardContent>
    </Card>
  )
}