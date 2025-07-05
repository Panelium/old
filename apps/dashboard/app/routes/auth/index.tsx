import { useState } from "react";
import { Card } from "~/components/ui/card";
import { AnimatePresence, motion } from "framer-motion";
import LoginPage from "./login";
import SignupPage from "./signup";

export default () => {
  const [signup, setSignup] = useState(false);

  return (
    <Card className="w-sm min-h-128 overflow-hidden flex items-center justify-center">
      <AnimatePresence mode="wait" initial={false}>
        {signup ? (
          <motion.div
            key="signup"
            initial={{ opacity: 0, x: 40 }}
            animate={{ opacity: 1, x: 0 }}
            exit={{ opacity: 0, x: -40 }}
            transition={{ duration: 0.35, ease: "easeInOut" }}
            className="w-full h-full flex items-center justify-center"
          >
            <SignupPage onLogin={() => setSignup(false)} />
          </motion.div>
        ) : (
          <motion.div
            key="login"
            initial={{ opacity: 0, x: -40 }}
            animate={{ opacity: 1, x: 0 }}
            exit={{ opacity: 0, x: 40 }}
            transition={{ duration: 0.35, ease: "easeInOut" }}
            className="w-full h-full flex items-center justify-center"
          >
            <LoginPage onSignUp={() => setSignup(true)} />
          </motion.div>
        )}
      </AnimatePresence>
    </Card>
  );
};
