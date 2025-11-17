import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { api } from "@/lib/api";
import { setLastLoginMethod } from "@/lib/last-login";
import { cn, sleep } from "@/lib/utils";
import { Link, useNavigate, useRouter } from "@tanstack/react-router";
import { motion } from "framer-motion";
import { Mail } from "lucide-react";
import { useForm } from "react-hook-form";
import { toast } from "sonner";
import Spinner from "../ui/spinner";
import { GoogleSVG } from "./google-svg";
import { useAuth } from "@/context/auth-context";
import { useMutation } from "@tanstack/react-query";

interface RegisterFormData {
  name: string;
  email: string;
  password: string;
}

export function RegisterForm({
  className,
  ...props
}: React.ComponentProps<"div">) {
  const navigate = useNavigate();
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<RegisterFormData>();

  const { setAuthState } = useAuth();
  const router = useRouter();

  const registerMutation = useMutation({
    mutationFn: async (data: RegisterFormData) => {
      return await api.post("/register", {
        username: data.name,
        email: data.email,
        password: data.password,
      });
    },
    onSuccess: async (response, data) => {
      setLastLoginMethod("email", data.email);
      setAuthState({ isAuthenticated: true });
      toast.success(response.data.message || "Account created successfully!");
      await sleep(1);
      navigate({ to: "/admin/home" });
    },
    onError: () => {
      toast.warning("Failed to create account");
    },
  });

  const handleRegister = async (data: RegisterFormData) => {
    registerMutation.mutate(data);
    router.invalidate();
  };

  const googleSignInMutation = useMutation({
    mutationFn: async () => {
      setLastLoginMethod("google");
      window.location.href = `/api/auth/google`;
    },
    onError: () => {
      toast.warning("Failed to sign in with Google");
      setAuthState({ isAuthenticated: false });
    },
  });

  const handleGoogleSignIn = (event: React.FormEvent) => {
    event.preventDefault();
    googleSignInMutation.mutate();
  };
  return (
    <div className={cn("flex flex-col gap-6", className)} {...props}>
      <Card>
        <CardHeader className="text-center">
          <CardTitle className="text-3xl instrument-serif-regular">
            Create an account
          </CardTitle>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit(handleRegister)}>
            <div className="grid gap-6">
              <div className="flex flex-col gap-4">
                <Button
                  onClick={handleGoogleSignIn}
                  variant="outline"
                  className="w-full"
                  type="button"
                >
                  <GoogleSVG />
                  Continue with Google
                </Button>
              </div>
              <div className="after:border-border relative text-center text-sm after:absolute after:inset-0 after:top-1/2 after:z-0 after:flex after:items-center after:border-t">
                <span className="bg-card text-muted-foreground relative z-10 px-2">
                  Or continue with
                </span>
              </div>
              <div className="grid gap-6">
                <div className="grid gap-3">
                  <Label htmlFor="name">Full Name</Label>
                  <Input
                    id="name"
                    type="text"
                    className="sm:text-xs text-sm"
                    placeholder="Enter your name"
                    {...register("name", {
                      required: "Your full name is required",
                      minLength: {
                        value: 2,
                        message: "Name must be at least 2 characters",
                      },
                    })}
                  />
                  {errors.name && (
                    <p className="text-sm text-red-500">
                      {errors.name.message}
                    </p>
                  )}
                </div>
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
                  <Label htmlFor="password">Password</Label>
                  <Input
                    id="password"
                    placeholder="Password"
                    type="password"
                    className="sm:text-xs text-sm"
                    {...register("password", {
                      required: "Password is required",
                      minLength: {
                        value: 8,
                        message: "Password must be at least 8 characters",
                      },
                      pattern: {
                        value: /^(?=.*[!@#$%^&*(),.?":{}|<>])/,
                        message:
                          "Password must contain at least one special character",
                      },
                    })}
                  />
                  {errors.password && (
                    <p className="text-sm text-red-500">
                      {errors.password.message}
                    </p>
                  )}
                </div>
                <motion.div
                  whileHover={{ scale: 1.02 }}
                  whileTap={{ scale: 0.98 }}
                >
                  <Button
                    type="submit"
                    className="w-full"
                    disabled={registerMutation.isPending}
                  >
                    {registerMutation.isPending ? (
                      <Spinner className="text-white" />
                    ) : (
                      <>
                        <Mail className="h-4 w-4" />
                        Continue with email
                      </>
                    )}
                  </Button>
                </motion.div>
              </div>
              <div className="text-center text-sm">
                Already have an account?{" "}
                <Link
                  to="/login"
                  className="underline font-semibold text-primary underline-offset-4"
                >
                  Log in
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