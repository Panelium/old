import React from "react";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { useNavigate } from "react-router-dom";
import { z } from "zod";
import { CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "~/components/ui/card";
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "~/components/ui/form";
import { Input } from "~/components/ui/input";
import { Button } from "~/components/ui/button";
import { useSession } from "~/providers/SessionProvider";
import { getAuthClient } from "~/lib/api-clients";
import { Turnstile } from "@marsidev/react-turnstile";
import { getConfig } from "~/lib/config";

const signupFormSchema = z.object({
  username: z.string().regex(/^[a-zA-Z0-9][a-zA-Z0-9_.-]{1,30}[a-zA-Z0-9]$/, {
    message:
      "Username must be 3-32 characters long, start and end with an alphanumeric character, and can contain letters, numbers, underscores, hyphens, and periods.",
  }),
  email: z.string().email({ message: "Please enter a valid email address" }),
  password: z
    .string()
    .min(16, {
      message: "Password must be at least 16 characters",
    })
    .max(384, {
      message: "Password must be at most 384 characters",
    }),
});

type SignupFormValues = z.infer<typeof signupFormSchema>;

export default function SignupPage({ onLogin }: { onLogin: () => void }) {
  const session = useSession();
  const navigate = useNavigate();

  const defaultValues: SignupFormValues = {
    username: "",
    email: "",
    password: "",
  };

  const form = useForm<SignupFormValues>({
    resolver: zodResolver(signupFormSchema),
    defaultValues,
  });

  const [turnstileStatus, setTurnstileStatus] = React.useState<"loaded" | "success" | "expired" | "error">("loaded");
  const [turnstileToken, setTurnstileToken] = React.useState<string>("");

  async function onSubmit(data: SignupFormValues) {
    const response = await (
      await getAuthClient()
    )
      .register({
        email: data.email,
        password: data.password,
        username: data.username,
        turnstileToken: turnstileToken,
      })
      .catch((error) => {
        console.error("Signup error:", error);
      });

    if (!response || !response.success) {
      form.setError("root", {
        message: "Signup failed. Please check your credentials.",
      });
      return;
    }

    session.setAuthenticated(true);
    navigate("/", { replace: true });
  }

  const [turnstileSiteKey, setTurnstileSiteKey] = React.useState<string>("");
  React.useEffect(() => {
    getConfig()
      .then((config) => {
        setTurnstileSiteKey(config.TURNSTILE_SITE_KEY);
      })
      .catch((error) => {
        console.error("Failed to load configuration:", error);
      });
  }, []);

  return (
    <div className="w-full max-w-sm truncate flex flex-col flex-1/2 gap-6 p-1">
      <CardHeader>
        <CardTitle className="text-2xl text-center no-select">Create an account</CardTitle>
        <CardDescription className="text-center no-select">
          Enter your credentials to access your account
        </CardDescription>
      </CardHeader>
      <CardContent>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
            <FormField
              control={form.control}
              name="username"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Username</FormLabel>
                  <FormControl>
                    <Input placeholder="username" className="no-select" {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="email"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Email</FormLabel>
                  <FormControl>
                    <Input placeholder="email@example.com" className="no-select" {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="password"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Password</FormLabel>
                  <FormControl>
                    <Input type="password" placeholder="••••••••" className="no-select" {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <Turnstile
              siteKey={turnstileSiteKey}
              onSuccess={(token) => {
                setTurnstileToken(token);
                setTurnstileStatus("success");
              }}
              onExpire={(token) => {
                setTurnstileToken("");
                setTurnstileStatus("expired");
              }}
              onError={(error) => {
                setTurnstileToken("");
                setTurnstileStatus("error");
                console.error("Turnstile error:", error);
                form.setError("root", {
                  message: "Turnstile verification failed. Please try again.",
                });
              }}
              options={{
                size: "invisible",
              }}
            />
            <Button type="submit" className="w-full no-select" disabled={turnstileStatus !== "success"}>
              Sign Up
            </Button>
            {form.formState.errors.root && (
              <p className="text-sm font-medium text-destructive text-center">{form.formState.errors.root.message}</p>
            )}
          </form>
        </Form>
      </CardContent>
      <CardFooter className="flex flex-col items-center">
        <p className="text-sm text-center text-muted-foreground no-select">
          Already have an account?{" "}
          <a className="cursor-pointer text-primary underline-offset-4 hover:underline" onClick={onLogin}>
            Return to login
          </a>
        </p>
      </CardFooter>
    </div>
  );
}
