import Link from "next/link";

type LegalSection = {
  title: string;
  body: string[];
  items?: string[];
};

type LegalPageProps = {
  title: string;
  description: string;
  updatedAt: string;
  sections: LegalSection[];
};

export function LegalPage({
  title,
  description,
  updatedAt,
  sections,
}: LegalPageProps) {
  return (
    <main className="lp-scope lp-legal">
      <div className="lp-legal-shell">
        <header className="lp-legal-top">
          <Link href="/" className="lp-legal-logo" aria-label="Entré トップへ">
            Entré
          </Link>
          <Link href="/login" className="lp-legal-back">
            ログインへ
          </Link>
        </header>

        <section className="lp-legal-hero">
          <p className="lp-simple-kicker">ENTRÉ BETA</p>
          <h1>{title}</h1>
          <p>{description}</p>
          <div className="lp-legal-updated">Last updated: {updatedAt}</div>
        </section>

        <div className="lp-legal-body">
          {sections.map((section) => (
            <section className="lp-legal-section" key={section.title}>
              <h2>{section.title}</h2>
              <div className="lp-legal-content">
                {section.body.map((paragraph) => (
                  <p key={paragraph}>{paragraph}</p>
                ))}
                {section.items ? (
                  <ul>
                    {section.items.map((item) => (
                      <li key={item}>{item}</li>
                    ))}
                  </ul>
                ) : null}
              </div>
            </section>
          ))}
        </div>
      </div>
    </main>
  );
}
