echo "Create account token with the default expiry"
vault write -force "${ENGINE_PATH}"/account/repo-reporting
echo "Create account token with a specific expiry"
vault write "${ENGINE_PATH}"/account/repo-reporting ttl=10s
vault write "${ENGINE_PATH}"/account/repo-reporting ttl=60s
echo "Create account token - capped expiry"
vault write "${ENGINE_PATH}"/account/repo-reporting ttl=50d

echo "Create project token with the default expiry"
vault write -force "${ENGINE_PATH}"/project/unprotected-cell-monitor/role/dev-role
echo "Create project token with a specific expiry"
vault write "${ENGINE_PATH}"/project/unprotected-cell-monitor/role/dev-role ttl=10s
vault write "${ENGINE_PATH}"/project/unprotected-cell-monitor/role/dev-role ttl=60s
echo "Create project token - capped expiry"
vault write "${ENGINE_PATH}"/project/unprotected-cell-monitor/role/dev-role ttl=2h