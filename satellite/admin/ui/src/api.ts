import type { Operation } from "./ui-generator";

interface API {
  authToken: string;
  operations: {
    [key: string]: [Operation];
  };
}

const api: API = {
  authToken: "",
  operations: {
    user: [
      {
        name: "create",
        desc: "Create a new user",
        params: [
          ["email", { type: "text", required: true }],
          ["full name", { type: "text", required: false }],
          ["password", { type: "password", required: true }],
          [
            "kind",
            {
              multiple: false,
              required: true,
              options: [
                { text: "", value: "" },
                { text: "personal", value: 1 },
                { text: "business", value: 2 },
              ],
            },
          ],
        ],
        func: async (
          email: string,
          fullName: string,
          password: string
        ): Promise<object> => {
          window.alert(`${email}, ${fullName}, ${password}`);
          return {};
        },
      },
    ],
  },
};

export default api;
