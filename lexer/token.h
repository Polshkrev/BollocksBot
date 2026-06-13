#ifndef TOKEN_H_
#define TOKEN_H_

#define STRING_VIEW_IMPLEMENTATION
#include "../build/Kada/lib/c/collections/string_view.h" // string_t, string_new

/**
 * @brief Representation of a specific and finite kind of tokens.
 */
typedef enum
{
    TokenEnd, // The end of the token stream.
    TokenColon, // `:`
    TokenWord, // A keyword within the token stream.
    TokenBang, // `!`
    TokenHash, // `#`
    TokenSingleQuote, // `''`
    TokenDoubleQuote, // `""`
    TokenUnknown, // An unknown token within the stream.
} TokenKind;

/**
 * @brief Representation of a token with a kind and a view.
 */
typedef struct
{
    TokenKind kind;
    string_t view;
} token_t;

#ifndef MAX_TOKENS
/** @brief Maximum allowable token count. This can be changed. */
#define MAX_TOKENS 256
#endif // MAX_TOKENS

/**
 * @brief Representation of a stream of tokens.
 */
typedef struct
{
    token_t tokens[MAX_TOKENS];
} token_array_t;

/**
 * @brief Obtain the string representation of a given token kind.
 * @param kind Kind of token to represent as a string.
 * @returns A string representation of the given token kind.
 */
const char *token_kind_to_string(TokenKind kind);

/**
 * @brief Construct a new token.
 * @returns A new token.
 */
token_t token_init(TokenKind kind, const char *data, size_t length);

#endif // TOKEN_H_

#ifdef TOKEN_IMPLEMENTATION

/**
 * @brief Obtain the string representation of a given token kind.
 * @param kind Kind of token to represent as a string.
 * @returns A string representation of the given token kind.
 */
const char *token_kind_to_string(TokenKind kind)
{
    switch (kind)
    {
        case TokenEnd:
        {
            return "TokenEnd";
        } break;
        case TokenColon:
        {
            return "TokenColon";
        } break;
        case TokenWord:
        {
            return "TokenWord";
        } break;
        case TokenBang:
        {
            return "TokenBang";
        } break;
        case TokenHash:
        {
            return "TokenHash";
        } break;
        case TokenSingleQuote:
        {
            return "TokenSingleQuote";
        } break;
        case TokenDoubleQuote:
        {
            return "TokenDoubleQuote";
        } break;
        case TokenUnknown:
        {
            return "TokenUnknown";
        } break;
        default:
        {
            return NULL;
        } break;
    }
}

/**
 * @brief Construct a new token.
 * @returns A new token.
 */
token_t token_init(TokenKind kind, const char *data, size_t length)
{
    return (token_t)
    {
        .kind = kind,
        .view = string_new(data, length)
    };
}

#endif // TOKEN_IMPLEMENTATION