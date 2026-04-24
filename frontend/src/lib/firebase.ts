// Firebase クライアント SDK の初期化。
//
// 重要: initializeApp を import 時（モジュール読み込み時）に呼ぶと
// Next.js のプリレンダリング時に env vars が無くて auth/invalid-api-key で落ちる。
// そのため lazy init にし、実際に getFirebaseAuth() が呼ばれたタイミングまで遅延する。
import { getApps, getApp, initializeApp, type FirebaseApp } from "firebase/app";
import { getAuth, GoogleAuthProvider, type Auth } from "firebase/auth";

const firebaseConfig = {
  apiKey: process.env.NEXT_PUBLIC_FIREBASE_API_KEY,
  authDomain: process.env.NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN,
  projectId: process.env.NEXT_PUBLIC_FIREBASE_PROJECT_ID,
  storageBucket: process.env.NEXT_PUBLIC_FIREBASE_STORAGE_BUCKET,
  messagingSenderId: process.env.NEXT_PUBLIC_FIREBASE_MESSAGING_SENDER_ID,
  appId: process.env.NEXT_PUBLIC_FIREBASE_APP_ID,
};

let cachedApp: FirebaseApp | null = null;
let cachedAuth: Auth | null = null;

export function getFirebaseAuth(): Auth {
  if (cachedAuth) return cachedAuth;
  cachedApp = getApps().length ? getApp() : initializeApp(firebaseConfig);
  cachedAuth = getAuth(cachedApp);
  return cachedAuth;
}

// GoogleAuthProvider は initializeApp 不要なモジュールトップでの生成でも安全
export const googleProvider = new GoogleAuthProvider();
