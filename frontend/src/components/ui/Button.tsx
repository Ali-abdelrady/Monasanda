import type { ButtonHTMLAttributes } from "react";

type Variant = "primary" | "secondary" | "danger" | "soft";

interface Props extends Omit<ButtonHTMLAttributes<HTMLButtonElement>, "type"> {
  variant?: Variant;
  type?: "button" | "submit";
}

export function Button({
  variant = "primary",
  className = "",
  type = "button",
  ...rest
}: Props) {
  return (
    <button
      type={type}
      className={`btn btn-${variant} ${className}`}
      {...rest}
    />
  );
}
