import type { Metadata } from "next";
import { LegalPage } from "@/components/legal/LegalPage";

export const metadata: Metadata = {
  title: "プライバシーポリシー | Entré",
  description: "Entré のプライバシーポリシーです。",
};

const sections = [
  {
    title: "1. 取得する情報",
    body: ["Entré は、サービス提供に必要な範囲で次の情報を取得します。"],
    items: [
      "Google ログインで取得されるユーザーID、メールアドレス、表示名",
      "応募先、企業、タスク、選考ステータス、StageHistory、保存した求人URLなど利用者が登録した情報",
      "アクセス日時、リクエスト情報、エラーログなどの技術的なログ",
    ],
  },
  {
    title: "2. 利用目的",
    body: ["取得した情報は、次の目的で利用します。"],
    items: [
      "本人確認、ログイン状態の維持、利用者ごとのデータ管理",
      "応募管理、企業管理、タスク管理、選考ステータス管理など本サービスの提供",
      "障害調査、不正利用防止、セキュリティ向上",
      "機能改善、品質改善、利用状況の把握",
    ],
  },
  {
    title: "3. Cookie",
    body: [
      "本サービスは、ログイン状態を維持するために httpOnly のセッション Cookie を利用します。この Cookie は、ブラウザ上の JavaScript から直接読み取れない形で扱います。",
    ],
  },
  {
    title: "4. 外部サービス",
    body: [
      "本サービスは、認証、ホスティング、API実行、データベース、ログ監視、メール送信などのために外部のクラウドサービスを利用する場合があります。",
    ],
    items: [
      "Firebase Authentication",
      "Vercel",
      "Google Cloud",
      "Supabase",
      "Resend",
    ],
  },
  {
    title: "5. 第三者提供",
    body: [
      "法令に基づく場合、利用者の同意がある場合、またはサービス提供に必要な委託先へ必要な範囲で取り扱いを委託する場合を除き、取得した個人情報を第三者に販売または提供しません。",
    ],
  },
  {
    title: "6. 安全管理",
    body: [
      "運営者は、認証、通信の暗号化、アクセス制御、ログ監視など、取得した情報の漏えい、滅失、改ざんを防ぐために必要な安全管理措置を講じます。",
    ],
  },
  {
    title: "7. 開示・削除",
    body: [
      "利用者から自身の情報の開示、訂正、削除、利用停止等の申し出があった場合、本人確認のうえ、法令に従って合理的な範囲で対応します。",
    ],
  },
  {
    title: "8. 改定",
    body: [
      "本ポリシーは、サービス内容や利用する外部サービスの変更に合わせて改定することがあります。重要な変更がある場合は、本サービス上または適切な方法で告知します。",
    ],
  },
];

export default function PrivacyPage() {
  return (
    <LegalPage
      title="プライバシーポリシー"
      description="Entré が取得する情報、利用目的、Cookie、外部サービスの扱いをまとめています。"
      updatedAt="2026-06-09"
      sections={sections}
    />
  );
}
