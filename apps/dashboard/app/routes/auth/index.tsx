import {
  ImperativePanelHandle,
  Panel,
  PanelGroup,
} from "react-resizable-panels";
import LoginPage from "./login";
import SignupPage from "./signup";
import { useRef, useState } from "react";
import { Card } from "~/components/ui/card";

export default () => {
  const [signup, setSignup] = useState(false);

  const loginRef = useRef<ImperativePanelHandle>(null);
  const signupRef = useRef<ImperativePanelHandle>(null);

  return (
    <Card className="max-w-sm max-h-130 h-1000">
      <PanelGroup direction="horizontal" className="items-center">
        <Panel
          defaultSize={100}
          className="transition-all duration-300 ease-in-out"
          ref={loginRef}
        >
          <LoginPage
            onSignUp={() => {
              setSignup(true);
              loginRef.current?.resize(0);
              signupRef.current?.resize(100);
            }}
          />
        </Panel>
        <Panel
          defaultSize={0}
          className="transition-all duration-300  ease-in-out"
          ref={signupRef}
        >
          <SignupPage
            onLogin={() => {
              setSignup(false);
              loginRef.current?.resize(100);
              signupRef.current?.resize(0);
            }}
          />
        </Panel>
      </PanelGroup>
    </Card>
  );
}
