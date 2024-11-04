"use client";

import { Suspense } from "react";
import VerifyUserForm from "@/app/components/veify_user/activation";

export default function VerifyUserPage() {
  return (
    <Suspense fallback={<LoadingState />}>
      <VerifyUserForm />
    </Suspense>
  );
}

function LoadingState() {
  return (
    <div className="min-h-screen flex items-center justify-center">
      <div className="text-center text-indigo-400">Loading...</div>
    </div>
  );
}

