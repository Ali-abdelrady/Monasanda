import { useState, type FormEvent } from "react";
import { ArrowRight, CircleAlert, ShieldCheck, UserRound } from "lucide-react";
import { Button } from "@/components/ui/Button";
import { useToast } from "@/context/ToastContext";

export function LoginPage({
  onLogin,
}: {
  onLogin: (email: string, password: string) => Promise<void>;
}) {
  const { notify } = useToast();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const submit = async (e: FormEvent) => {
    e.preventDefault();
    setError("");
    setLoading(true);
    try {
      await onLogin(email, password);
      notify("تم تسجيل الدخول بنجاح");
    } catch (err) {
      setError(
        err instanceof Error
          ? err.message
          : "تعذر تسجيل الدخول، راجع البيانات وحاول مرة أخرى",
      );
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="login-layout">
      <LoginArt />
      <div className="login-side">
        <form className="login-form" onSubmit={submit}>
          <div className="mobile-brand">
            <div className="brand-mark">
              <ShieldCheck size={24} />
            </div>
            <span>تمامات</span>
          </div>
          <div className="eyebrow">أهلاً بعودتك</div>
          <h2>تسجيل الدخول</h2>
          <p className="muted">سجّل الدخول للوصول إلى دفعات الإجازات وإدارة الحضور.</p>

          <label className="field-label" htmlFor="email">
            البريد الإلكتروني
          </label>
          <div className="input-wrap">
            <UserRound size={18} />
            <input
              id="email"
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="name@example.com"
              autoComplete="username"
              required
              dir="ltr"
            />
          </div>

          <label className="field-label" htmlFor="password">
            كلمة المرور
          </label>
          <div className="input-wrap">
            <ShieldCheck size={18} />
            <input
              id="password"
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="أدخل كلمة المرور"
              autoComplete="current-password"
              required
            />
          </div>

          {error && (
            <div className="inline-error">
              <CircleAlert size={16} />
              {error}
            </div>
          )}
          <Button type="submit" disabled={loading} className="login-submit">
            {loading ? (
              <>
                <span className="spinner small" /> جاري تسجيل الدخول
              </>
            ) : (
              <>
                دخول إلى النظام <ArrowRight size={18} />
              </>
            )}
          </Button>
          <div className="login-security">
            <ShieldCheck size={16} />
            <span>جلسة دخول آمنة ومحمية</span>
          </div>
        </form>
        <div className="login-copyright">نظام داخلي لإدارة ومتابعة الحضور</div>
      </div>
    </div>
  );
}

function LoginArt() {
  return (
    <div className="login-art">
      <div className="art-orbit orbit-one" />
      <div className="art-orbit orbit-two" />
      <div className="art-content">
        <div className="brand-mark brand-large">
          <ShieldCheck size={34} />
        </div>
        <div className="eyebrow light-eyebrow">نظام الإدارة والمتابعة</div>
        <h1>
          التمام..
          <br />
          <span>بشكل أذكى.</span>
        </h1>
        <p>
          منصة موحدة لإدارة دفعات الإجازات، وتسجيل الحضور، ومتابعة عودة الأفراد لحظة
          بلحظة.
        </p>
        <div className="art-stat-row">
          <div>
            <strong>متابعة دقيقة</strong>
            <span>كل فرد في مكانه الصحيح</span>
          </div>
          <div>
            <strong>تحديث مستمر</strong>
            <span>بيانات متزامنة بين الأجهزة</span>
          </div>
        </div>
      </div>
      <div className="art-footer">
        تمامات <span>•</span> إدارة عودة الإجازات
      </div>
    </div>
  );
}
