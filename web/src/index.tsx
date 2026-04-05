import { render } from "solid-js/web";
import { Router, Route } from "@solidjs/router";
import Layout from "./components/Layout.tsx";
import LabListPage from "./pages/LabListPage.tsx";
import LabPage from "./pages/LabPage.tsx";
import NotFoundPage from "./pages/NotFoundPage.tsx";

const root = document.getElementById("app");
if (root) {
  render(
    () => (
      <Router root={Layout}>
        <Route path="/" component={LabListPage} />
        <Route path="/labs/:id" component={LabPage} />
        <Route path="*" component={NotFoundPage} />
      </Router>
    ),
    root,
  );
}
