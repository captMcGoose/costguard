package pricing

var EC2Pricing = map[string]float64{
    "t3.micro": 8.0,
    "t3.small": 15.0,
    "t3.medium": 30.0,
}

var RDSPricing = map[string]float64{
    "db.t3.micro": 15.0,
    "db.t3.small": 30.0,
    "db.m6i.large": 180.0,
}

const NATGatewayMonthly = 32.0
const EBSPricePerGB = 0.08
