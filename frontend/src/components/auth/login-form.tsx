import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { api } from "@/lib/api";
import { getLastLoginMethod, setLastLoginMethod } from "@/lib/last-login"
import { cn, sleep } from "@/lib/utils";
import { Link, useNavigate, useRouter } from "@tanstack/react-router";
import { motion } from "framer-motion";
import { Mail } from "lucide-react";
import { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { toast } from "sonner";
import Spinner from "../ui/spinner";
import { GoogleSVG } from "./google-svg";
import { useAuth } from "@/context/auth-context";
import { useMutation } from "@tanstack/react-query";

interface LoginFormData {
  email: string;
  password: string;
}

export function LoginForm({
  className,
  ...props
}: React.ComponentProps<"div">) {
  const navigate = useNavigate();
  const {
    register,
    handleSubmit,
    formState: { errors },
    setValue,
  } = useForm<LoginFormData>();

  const [lastLogin, setLastLogin] = useState<ReturnType<typeof getLastLoginMethod>>(null);
  const router = useRouter();
  const { setAuthState } = useAuth();


  useEffect(() => {
    const lastLoginInfo = getLastLoginMethod();
    setLastLogin(lastLoginInfo);

    if (lastLoginInfo?.method === 'email' && lastLoginInfo.email) {
      setValue('email', lastLoginInfo.email);
    }
  }, [setValue]);


  const loginMutation = useMutation({
    mutationFn: async (data: LoginFormData) => {
      return api.post("/login", {
        email: data.email,
        password: data.password,
      });
    },
    onSuccess: async (response, variables) => {
      setLastLoginMethod("email", variables.email);
      setAuthState({ isAuthenticated: true, isLoading: false });
      toast.success(response.data.message || "Login successful!");
      await sleep(1)
      navigate({ to: "/admin/home" });
    },
    onError: () => {
      toast.warning("Invalid email or password");
    },
  });

  const handleEmailSignIn = async (data: LoginFormData) => {
    loginMutation.mutate(data);
    router.invalidate()
  };

  const googleLoginMutation = useMutation({
    mutationFn: async () => {
      setLastLoginMethod("google");
      return { redirectUrl: `/api/auth/google` };
    },
    onSuccess: (data) => {
      window.location.href = data.redirectUrl;
    },
    onError: () => {
      toast.warning("Failed to sign in with Google");
    },
  });

  const handleGoogleSignIn = async (event: React.FormEvent) => {
    event.preventDefault();
    googleLoginMutation.mutate();
  };

  return (
    <div className={cn("flex flex-col gap-6", className)} {...props}>
      <Card>
        <CardHeader className="text-center">
          <CardTitle className="text-3xl instrument-serif-regular">
            Welcome back
          </CardTitle>

        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit(handleEmailSignIn)}>
            <div className="grid gap-6">
              <div className="flex flex-col gap-4">
                <Button
                  onClick={handleGoogleSignIn}
                  variant="outline"
                  className={cn(
                    "w-full relative",
                    lastLogin?.method === 'google' && "ring-2 ring-primary/20 bg-primary/5"
                  )}
                  type="button"
                >
                  <GoogleSVG />
                  Continue with Google
                  {lastLogin?.method === 'google' && (
                    <div className="absolute -top-1 -right-1 flex items-center gap-1 bg-primary text-primary-foreground text-xs px-2 py-0.5 rounded-full">
                      Last used
                    </div>
                  )}
                </Button>
              </div>
              <div className="after:border-border relative text-center text-sm after:absolute after:inset-0 after:top-1/2 after:z-0 after:flex after:items-center after:border-t">
                <span className="bg-card text-muted-foreground relative z-10 px-2">
                  Or continue with
                </span>
              </div>
              <div className="grid gap-6">
                <div className="grid gap-3">
                  <Label htmlFor="email">Email</Label>
                  <Input
                    id="email"
                    type="email"
                    className="sm:text-xs text-sm"
                    placeholder="jane@example.com"
                    {...register("email", {
                      required: "Email is required",
                      pattern: {
                        value: /^\S+@\S+$/i,
                        message: "Invalid email address",
                      },
                    })}
                  />
                  {errors.email && (
                    <p className="text-sm text-red-500">
                      {errors.email.message}
                    </p>
                  )}
                </div>
                <div className="grid gap-3">
                  <div className="flex items-center">
                    <Label htmlFor="password">Password</Label>
                    <Link
                      to="/forgot-password"
                      className="ml-auto text-sm text-muted-foreground underline-offset-4 hover:underline"
                    >
                      Forgot your password?
                    </Link>
                  </div>
                  <Input
                    placeholder="Password"
                    id="password"
                    type="password"
                    className="sm:text-xs text-sm"
                    {...register("password", {
                      required: "Password is required",
                      minLength: {
                        value: 6,
                        message: "Password must be at least 6 characters",
                      },
                      pattern: {
                        value: /^(?=.*[!@#$%^&*(),.?":{}|<>])/,
                        message:
                          "Password must contain at least one special character",
                      },
                    })}
                  />
                  {errors.password && (
                    <p className="text-sm  text-red-500">
                      {errors.password.message}
                    </p>
                  )}
                </div>
                <motion.div
                  whileHover={{ scale: 1.02 }}
                  whileTap={{ scale: 0.98 }}
                  className="relative"
                >
                  <Button
                    type="submit"
                    className={cn(
                      "w-full relative",
                      lastLogin?.method === 'email' && "ring-2 ring-primary/20"
                    )}
                    disabled={loginMutation.isPending}
                  >
                    {loginMutation.isPending ? (
                      <Spinner className="text-white" />
                    ) : (
                      <>
                        <Mail className="h-4 w-4" />
                        Continue with email
                      </>
                    )}
                    {lastLogin?.method === 'email' && (
                      <div className="absolute shadow-sm -top-1 -right-1 flex items-center gap-1 bg-white text-primary text-xs px-2 py-0.5 rounded-full">
                        Last used
                      </div>
                    )}
                  </Button>

                </motion.div>
              </div>
              <div className="text-center text-sm">
                Don&apos;t have an account?{" "}
                <Link
                  to="/register"
                  className="underline font-semibold text-primary underline-offset-4"
                >
                  Sign up
                </Link>
              </div>
            </div>
          </form>
        </CardContent>
      </Card>
      <div className="text-muted-foreground *:[a]:hover:text-primary text-center text-xs text-balance *:[a]:underline *:[a]:underline-offset-4">
        By continuing, you agree to our{" "}
        <Link to="/terms">Terms of Service</Link> and{" "}
        <Link to="/privacy">Privacy Policy</Link>.
      </div>
    </div>
  );
}