echo "Create account token with the default expiry"
vault write -force "${ENGINE_PATH}"/account/wfe-ops
echo "Create account token with a specific expiry"
vault write "${ENGINE_PATH}"/account/wfe-ops ttl=10s
vault write "${ENGINE_PATH}"/account/wfe-ops ttl=60s
echo "Create account token - capped expiry"
vault write "${ENGINE_PATH}"/account/wfe-ops ttl=50d

echo "Create project token with the default expiry"
vault write -force "${ENGINE_PATH}"/project/wfecd-stg-unprotected/role/unprotected-role
echo "Create project token with a specific expiry"
vault write "${ENGINE_PATH}"/project/wfecd-stg-unprotected/role/unprotected-role ttl=10s
vault write "${ENGINE_PATH}"/project/wfecd-stg-unprotected/role/unprotected-role ttl=60s
echo "Create project token - capped expiry"
vault write "${ENGINE_PATH}"/project/wfecd-stg-unprotected/role/unprotected-role ttl=2h