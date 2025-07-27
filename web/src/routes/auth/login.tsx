import { createFileRoute } from '@tanstack/react-router'
import { LoginForm } from '@/components/LoginForm'

export const Route = createFileRoute('/auth/login')({
  component: LoginPage,
})

function LoginPage() {
  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md w-full space-y-8">
        <div className="text-center">
          <h2 className="mt-6 text-3xl font-extrabold text-gray-900">
            JobWatch Canada
          </h2>
          <p className="mt-2 text-sm text-gray-600">
            Sign in to access business demographics and reporting
          </p>
        </div>
        <LoginForm />
      </div>
    </div>
  )
}