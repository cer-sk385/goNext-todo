provider "azurerm" {
  features {}
  subscription_id = var.subscription_id
}

// 既存のリソースグループをdata sourceとして定義
data "azurerm_resource_group" "shin_private_rg" {
  name = var.resource_group_name
}

// AKSクラスターの作成
resource "azurerm_kubernetes_cluster" "aks" {
  name                = "example-aks"
  location            = data.azurerm_resource_group.shin_private_rg.location
  resource_group_name = data.azurerm_resource_group.shin_private_rg.name
  dns_prefix          = "exampleaks"

  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "Standard_D2_v2"
  }

  identity {
    type = "SystemAssigned"
  }

  tags = {
    Environment = "Development"
  }
}

// ACRの作成
resource "azurerm_container_registry" "acr" {
  name                = "skprivateacr"  // ACR名は一意である必要があります
  resource_group_name = data.azurerm_resource_group.shin_private_rg.name
  location            = data.azurerm_resource_group.shin_private_rg.location
  sku                = "Standard"
  admin_enabled      = false  // マネージドIDを使用するため、admin認証は無効化

  tags = {
    Environment = "Development"
  }
}

// AKSにACRへのPull権限を付与
resource "azurerm_role_assignment" "aks_acr" {
  principal_id                     = azurerm_kubernetes_cluster.aks.kubelet_identity[0].object_id
  role_definition_name            = "AcrPull"
  scope                           = azurerm_container_registry.acr.id
  skip_service_principal_aad_check = true
}

// AKSクラスターのクレデンシャル出力
output "kube_config" {
  value     = azurerm_kubernetes_cluster.aks.kube_config_raw
  sensitive = true
}

// ACRのログイン情報を出力
output "acr_login_server" {
  value = azurerm_container_registry.acr.login_server
}
