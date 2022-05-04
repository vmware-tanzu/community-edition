import React, { useReducer } from "react";
import { globalAppStateReducer, initialState, IAppState } from "./reducer";

interface IAppStateProps {
  state: IAppState;
  dispatch: ({ type }: { type: string }) => void;
}

export const GlobalAppState = React.createContext({} as IAppStateProps);

function init(initialState: IAppState){
  const stickyValue: string | null = window.localStorage.getItem("appState");
  return stickyValue !== null ? JSON.parse(stickyValue).value : initialState;
}

export function GlobalAppStateProvider(props: any) {
  const [state, dispatch] = useReducer(globalAppStateReducer, initialState, init);

  const value = { state, dispatch } as IAppStateProps;
  return <GlobalAppState.Provider value={value}>{props.children}</GlobalAppState.Provider>;
}