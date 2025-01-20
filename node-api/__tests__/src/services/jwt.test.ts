import { validateToken } from "@Services/jwt"
import { Request } from 'express';
import { Config } from '@Services/env'
import { createPublicKey } from "node:crypto";
import fs from 'fs'
import path from "path";

function readFromFile(relPath: string): string {
  return fs.readFileSync(path.resolve(__dirname, relPath)).toString()
}

jest.mock('@Services/env', () => ({
  environment: jest.fn(() => 
    new Promise((resolve) => {
      resolve({
        env: {},
        secEnv: {
          jwtKeyPublic: createPublicKey(readFromFile("./test-jwt-rsa256-public.pem")),
        }
      } as Config)
    }),
)}));

test('GIVEN token is valid THEN validateToken returns true', () => {
  const mockRequest: unknown = {
    headers: {
      cookie: `jwt_token=${readFromFile("./validToken.txt")}`
   }
  };
  validateToken(mockRequest as Request).then((value) => {
    expect(value).toBe(true);
  })
});

test('GIVEN token is missing THEN validateToken returns false', () => {
  const mockRequest: unknown = {
    headers: {
      cookie: ""
   }
  };
  validateToken(mockRequest as Request).then((value) => {
    expect(value).toBe(false);
  })
});

test('GIVEN token has invalid signature THEN validateToken returns false', () => {
  const mockRequest: unknown = {
    headers: {
      cookie: `jwt_token=${readFromFile("./validToken.txt")}badsignature`
   }
  };
  validateToken(mockRequest as Request).then((value) => {
    expect(value).toBe(false);
  })
});

test('GIVEN token has expired THEN validateToken returns false', () => {
  const mockRequest: unknown = {
    headers: {
      cookie: `jwt_token=${readFromFile("./expiredToken.txt")}`
   }
  };
  validateToken(mockRequest as Request).then((value) => {
    expect(value).toBe(false);
  })
});
