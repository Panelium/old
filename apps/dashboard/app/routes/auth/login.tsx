import React from "react";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { useNavigate } from "react-router-dom";
import { z } from "zod";
import {
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "~/components/ui/card";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "~/components/ui/form";
import { Input } from "~/components/ui/input";
import { Button } from "~/components/ui/button";
import { Checkbox } from "~/components/ui/checkbox";
import { useSession } from "~/providers/SessionProvider";
import { getAuthClient } from "~/lib/api-clients";

const loginFormSchema = z.object({
  email: z
    .string()
    .email({ message: "Please enter a valid username or email address" })
    .or(
      z.string().regex(/^[a-zA-Z0-9][a-zA-Z0-9_.-]{1,30}[a-zA-Z0-9]$/, {
        message: "Please enter a valid username or email address",
      })
    ),
  password: z
    .string()
    .min(16, {
      message: "Please enter a valid password",
    })
    .max(384, {
      message: "Please enter a valid password",
    }),
  rememberMe: z.boolean(),
});

type LoginFormValues = z.infer<typeof loginFormSchema>;

export default function LoginPage({ onSignUp }: { onSignUp?: () => void }) {
  const session = useSession();
  const navigate = useNavigate();

  const defaultValues: LoginFormValues = {
    email: "",
    password: "",
    rememberMe: false,
  };

  const form = useForm<LoginFormValues>({
    resolver: zodResolver(loginFormSchema),
    defaultValues,
  });

  async function onSubmit(data: LoginFormValues) {
    const response = await (
      await getAuthClient()
    )
      .login({
        username: data.email,
        password: data.password,
      })
      .catch((error) => {
        console.error("Login error:", error);
      });

    if (!response || !response.success) {
      form.setError("root", {
        message: "Login failed. Please check your credentials.",
      });
      return;
    }

    session.setAuthenticated(true);
    navigate("/", { replace: true });
  }

  return (
    <div className="w-full max-w-sm truncate flex flex-col flex-1/2 gap-6 p-1">
      <CardHeader>
        <CardTitle className="text-2xl text-center no-select">
          Sign In
        </CardTitle>
        <CardDescription className="text-center no-select">
          Enter your credentials to access your account
        </CardDescription>
      </CardHeader>
      <CardContent>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
            <FormField
              control={form.control}
              name="email"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Email or Username</FormLabel>
                  <FormControl>
                    <Input
                      placeholder="email@example.com"
                      className="no-select"
                      {...field}
                    />
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
                    <Input
                      type="password"
                      placeholder="••••••••"
                      className="no-select"
                      {...field}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="rememberMe"
              render={({ field }) => (
                <FormItem className="flex flex-row items-center space-x-3 space-y-0">
                  <FormControl>
                    <Checkbox
                      checked={field.value}
                      onCheckedChange={field.onChange}
                    />
                  </FormControl>
                  <div className="space-y-1 leading-none no-select">
                    <FormLabel>Remember me</FormLabel>
                  </div>
                </FormItem>
              )}
            />
            <Button type="submit" className="w-full no-select">
              Sign In
            </Button>
            {form.formState.errors.root && (
              <p className="text-sm font-medium text-destructive text-center">
                {form.formState.errors.root.message}
              </p>
            )}
          </form>
        </Form>
      </CardContent>
      <CardFooter className="flex flex-col items-center">
        <p className="text-sm text-center text-muted-foreground no-select">
          Don't have an account?{" "}
          <a
            className="cursor-pointer text-primary underline-offset-4 hover:underline"
            onClick={onSignUp}
          >
            Sign up
          </a>
        </p>
      </CardFooter>
    </div>
  );
}
